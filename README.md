# go-phylum

API Wrapper for Phylum's REST API

**Incomplete** This API wrapper does not implement all functionality yet

## Quickstart
```golang
package main

import "fmt"
import "github.com/peterjmorgan/go-phylum"

func main() {
	// Create Client using locally-installed Phylum CLI as source of oauth refresh token 
	client, err := phylum.NewClient(&phylum.ClientOptions{})
        if err != nil {
		fmt.Printf("Failed to create Phylum client: %v\n", err)
        }
	
	projects, err := client.ListProjects()
	if err != nil {
		fmt.Printf("Failed to list projects: %v\n", err)
	}
	_ = projects
	
	// Create client with user-supplied oauth refresh token
	client2, err2 := phylum.NewClient(&phylum.ClientOptions{
		Token: "user-supplied value of refresh token",
        })
	if err2 != nil {
		fmt.Printf("Failed to create Phylum client: %v\n", err2)
	}
	
	// Create a Phylum project to store analysis results
	project, err := client.CreateProject("myProject", &ProjectOpts{})
	if err != nil {
		fmt.Printf("Failed to create project: %v\n", err)
	}
	
	// Parse a local lockfile file using Phylum's API
	packages, err := client.ParseLockfile("poetry.lock")
	if err != nil {
		fmt.Printf("Failed to parse lockfile: %v\n", err)
	}
	
	// Submit packages to Phylum for analysis, returning a job identifier
	jobId, err := client.AnalyzeParsedPackages(*project.Ecosystem, project.Id, packages)
	if err != nil {
		fmt.Printf("Failed to analyze packages: %v\n", err)
	}
	
	// Request complete job results
	results, err := client.GetJobVerbose(jobId)
	if err != nil {
		fmt.Printf("Failed to request job results: %v\n", err)
	}
	
	// Print Phylum Project Score from analysis
	fmt.Printf("Phylum project score: %v\n", results.Score)
	
	// Iterate through analyzed packages printing the Phylum Package Score for each
	for _, package := range results.Packages {
		fmt.Printf("%v:%v@%v - %v\n", package.Type, package.Name, package.Version, package.PackageScore)
	}
}
```
