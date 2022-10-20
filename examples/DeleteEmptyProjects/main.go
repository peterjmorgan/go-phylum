package main

import (
	"fmt"

	"github.com/peterjmorgan/go-phylum"
)

func main() {
	p, err := phylum.NewClient(&phylum.ClientOptions{})
	if err != nil {
		fmt.Printf("Failed to create phylum client: %v\n", err)
		return
	}

	allProjects, err := p.GetAllProjects()
	if err != nil {
		fmt.Printf("Failed to GetAllProjects(): %v\n", err)
		return
	}
	fmt.Printf("Found %v Phylum Projects\n", len(allProjects))

	var count int
	for _,proj := range allProjects {
		if len(proj.Dependencies) == 0 {
			count++
			_, err := p.DeleteProject(proj.Id.String())
			if err != nil {
				fmt.Printf("Failed to DeleteProject with name %v: %v\n", proj.Name, err)
			}
		}
	}

	fmt.Printf("[DONE] Deleted %v/%v projects\n", count, len(allProjects))
}
