package main

import (
	"encoding/json"
	"fmt"
	"os"

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

	projectData, err := json.Marshal(allProjects)
	if err != nil {
		fmt.Printf("Failed to marshall JSON: %v\n", err)
		return
	}

	if err := os.WriteFile("allProjectData.json", projectData, 0644); err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
	}
}