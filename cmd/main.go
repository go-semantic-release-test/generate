package main

import (
	"github.com/go-semantic-release-test/generate"
	"log"
)

func main() {
	err := generate.Generate()

	if err != nil {
		log.Fatal(err)
	}
}
