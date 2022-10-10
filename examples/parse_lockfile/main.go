package main

import (
	"fmt"

	phylum "github.com/peterjmorgan/go-phylum"
)

func main() {
	p, err := phylum.NewClient(&phylum.ClientOptions{})
	if err != nil {
		fmt.Printf("failed to create PhylumClient: %v\n", err)
		return
	}

	packages, err := p.ParseLockfile("../../test_lockfiles/package-lock.json")
	if err != nil {
		fmt.Printf("failed to parse lockfile: %v\n", err)
		return
	}

	_ = packages
}
