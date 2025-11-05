// Package storage provides workflow storage and retrieval functionality.
//
// This package implements an in-memory storage system for workflows,
// allowing workflows to be saved, loaded, listed, and deleted by ID.
//
// # Usage
//
//	store := storage.NewInMemoryStore()
//
//	// Save a workflow
//	id, err := store.Save("my-workflow", workflowData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Load a workflow
//	workflow, err := store.Load(id)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// List all workflows
//	workflows := store.List()
//
// # Security Considerations
//
// The in-memory store is suitable for development and testing but should
// not be used in production without persistence. For production use,
// consider implementing a persistent storage backend.
package storage
