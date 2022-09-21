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
	
	// Create client with user-supplied oauth refresh token
	client2, err2 := phylum.NewClient(&phylum.ClientOptions{
		Token: "user-supplied value of refresh token",
        })
	if err2 != nil {
		fmt.Printf("Failed to create Phylum client: %v\n", err2)
	}
}
```