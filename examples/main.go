package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/istovpets/sqlset"
)

//go:embed queries
var queriesFS embed.FS

func main() {
	// Create a new SQLSet from the embedded filesystem.
	sqlSet, err := sqlset.New(queriesFS)
	if err != nil {
		log.Fatalf("Failed to create SQL set: %v", err)
	}

	// Use interfaces
	var queries sqlset.SQLQueriesProvider = sqlSet
	var sets sqlset.SQLSetsProvider = sqlSet

	// Get a specific query
	query, err := queries.Get("users", "GetUserByID")
	if err != nil {
		log.Fatalf("Failed to get query: %v", err)
	}
	fmt.Println("GetUserByID query (multi argument):", query)

	query, err = queries.Get("users.CreateUser")
	if err != nil {
		log.Fatalf("Failed to get query: %v", err)
	}
	fmt.Println("CreateUser query (dot notation):", query)

	query, err = queries.Get("CreateUser")
	if err != nil {
		log.Fatalf("Failed to get query: %v", err)
	}
	fmt.Println("CreateUser query (single argument):", query)

	// Or, panic if the query is not found
	query = queries.MustGet("users", "CreateUser")
	fmt.Println("CreateUser query (MustGet method):", query)

	fmt.Println("--------------------------------")
	// You can also retrieve metadata for all query sets
	metas := sets.GetSetsMetas()
	for _, meta := range metas {
		fmt.Printf("Set ID: %s, Name: %s, Description: %s\n", meta.ID, meta.Name, meta.Description)
	}

	// You can get a list of all query IDs in a specific set
	queryIDs, err := sets.GetQueryIDs("users")
	if err != nil {
		log.Fatalf("Failed to get query IDs: %v", err)
	}
	fmt.Println("Query IDs in 'users' set:", queryIDs)
}
