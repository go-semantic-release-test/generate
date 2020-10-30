package generate

import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-semantic-release-test/generate/scenarios"
	"os"
	"time"
)

type Generator struct {
	repo *git.Repository
}

func NewGenerator(repo *git.Repository) *Generator {
	return &Generator{repo: repo}
}

func CreateAuthor() *object.Signature {
	return &object.Signature{
		Name:  "John Doe",
		Email: "john@doe.org",
		When:  time.Now(),
	}
}

const defaultMessage = "Some Commit Message"

func (g Generator) Commit(change scenarios.ChangeType, msg string) (*plumbing.Hash, error) {
	w, err := g.repo.Worktree()

	if err != nil {
		return nil, err
	}

	if msg == "" {
		msg = defaultMessage
	}

	prefix := ""

	switch change {
	case scenarios.Feat:
		prefix = "feat: "
	case scenarios.Fix:
		prefix = "fix: "
	case scenarios.Breaking:
		prefix = "feat: BREAKING "
	case scenarios.Chore:
		prefix = "chore: "
	case scenarios.Other:
		// No Prefix
	default:
		return nil, fmt.Errorf("invalid change type: %s", change)
	}

	hash, err := w.Commit(fmt.Sprintf("%s%s", prefix, msg), &git.CommitOptions{
		Author: CreateAuthor(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return &hash, nil
}

func (g Generator) Run(scenarios []scenarios.Scenario) error {
	for _, sc := range scenarios {
		w, err := g.repo.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree: %w", err)
		}

		err = g.processScenario(w, sc)

		if err != nil {
			return err
		}
	}

	return nil
}

func (g Generator) processScenario(w *git.Worktree, scenario scenarios.Scenario) error {
	err := w.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/master",
	})

	if err != nil {
		return fmt.Errorf("failed to checkout master branch: %w", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(scenario.BranchName),
		Create: true,
	})

	if err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	for _, c := range scenario.Commits {
		hash, err := g.Commit(c.ChangeType, c.Message)

		if err != nil {
			return err
		}

		// Create tag
		if c.Tag != nil {
			_, err = g.repo.CreateTag(c.Tag.Name, *hash, nil)

			if err != nil {
				return fmt.Errorf("failed to create tag: %w", err)
			}
		}
	}

	return nil
}

func Generate() error {
	auth, err := ssh2.NewPublicKeys("git", []byte(os.Getenv("PRIVATE_KEY")), "")

	if err != nil {
		return fmt.Errorf("failed to create auth credentials: %w", err)
	}

	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:        "github.com:go-semantic-release-test/test.git",
		Auth:       auth,
		RemoteName: "github",
		Progress:   os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "gitlab",
		URLs: []string{"git@gitlab.com:go-semantic-release/test/test.git"},
	})

	if err != nil {
		return fmt.Errorf("failed to create gitlab remote")
	}

	g := NewGenerator(repo)

	err = g.Run(scenarios.Scenarios)

	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{
		Auth:       auth,
		RemoteName: "github",
		Progress:   os.Stdout,
		Force:      true,
	})

	if err != nil {
		return fmt.Errorf("failed to push to github")
	}

	err = repo.Push(&git.PushOptions{
		RemoteName: "gitlab",
		Auth:       auth,
		Progress:   os.Stdout,
		Force:      true,
	})

	if err != nil {
		return fmt.Errorf("failed to push to gitlab")
	}

	return nil
}
