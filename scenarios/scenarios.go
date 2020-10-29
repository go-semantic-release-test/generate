package scenarios

import "fmt"

type ChangeType string

const (
	Feat     ChangeType = "Feat"
	Fix      ChangeType = "Fix"
	Chore    ChangeType = "Chore"
	Breaking ChangeType = "Breaking"
	Other    ChangeType = "Other"
)

type Scenario struct {
	Description string
	BranchName  string
	Commits     []Commit
}

type Tag struct {
	Name string
}

type Commit struct {
	Message    string
	ChangeType ChangeType
	Tag        *Tag // Optional tag
}

func generateCommits(changeType ChangeType, n int) []Commit {
	i := 0
	var commits []Commit

	for i < n {
		commits = append(commits, Commit{
			Message:    fmt.Sprintf("Commit %d", i),
			ChangeType: changeType,
		})

		i += 1
	}

	return commits
}

func init() {
	var scenario1 = Scenario{
		Description: "No Commits",
		BranchName:  "scenario1",
	}

	var scenario2 = Scenario{
		Description: "1 Fix",
		BranchName:  "scenario2",
		Commits: []Commit{{
			ChangeType: Fix,
		}},
	}
	var scenario3 = Scenario{
		Description: "1 Fix",
		BranchName:  "scenario3",
		Commits: []Commit{{
			ChangeType: Fix,
		}},
	}

	var scenario4 = Scenario{
		Description: "Lots of Commits",
		BranchName:  "scenario4",
		Commits:     generateCommits(Fix, 10000),
	}

	var scenario5 = Scenario{
		Description: "External Tag",
		BranchName:  "scenario5",
		Commits: []Commit{{
			ChangeType: Fix,
		}, {
			ChangeType: Fix,
		}, {
			ChangeType: Feat,
		}, {
			ChangeType: Breaking,
			Tag:        &Tag{Name: "v1.0.0"},
		}, {
			ChangeType: Fix,
		}},
	}

	Scenarios = []Scenario{scenario1, scenario2, scenario3, scenario4, scenario5}
}

var Scenarios []Scenario
