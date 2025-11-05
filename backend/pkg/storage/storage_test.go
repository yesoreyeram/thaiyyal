package storage

import (
	"encoding/json"
	"testing"
)

func TestInMemoryStore_Save(t *testing.T) {
	store := NewInMemoryStore()
	
	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	
	tests := []struct {
		name        string
		workflowName string
		description string
		data        json.RawMessage
		wantErr     bool
	}{
		{
			name:        "Valid workflow",
			workflowName: "Test Workflow",
			description: "A test workflow",
			data:        data,
			wantErr:     false,
		},
		{
			name:        "Empty name",
			workflowName: "",
			description: "Description",
			data:        data,
			wantErr:     true,
		},
		{
			name:        "Empty data",
			workflowName: "Test",
			description: "Description",
			data:        json.RawMessage{},
			wantErr:     true,
		},
		{
			name:        "Invalid JSON data",
			workflowName: "Test",
			description: "Description",
			data:        json.RawMessage(`{invalid json`),
			wantErr:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := store.Save(tt.workflowName, tt.description, tt.data)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if id == "" {
				t.Error("Expected non-empty ID")
			}
		})
	}
}

func TestInMemoryStore_Load(t *testing.T) {
	store := NewInMemoryStore()
	
	data := json.RawMessage(`{"nodes": [{"id": "1"}], "edges": []}`)
	id, err := store.Save("Test Workflow", "Description", data)
	if err != nil {
		t.Fatalf("Failed to save workflow: %v", err)
	}
	
	t.Run("Load existing workflow", func(t *testing.T) {
		workflow, err := store.Load(id)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		
		if workflow.ID != id {
			t.Errorf("Expected ID %s, got %s", id, workflow.ID)
		}
		
		if workflow.Name != "Test Workflow" {
			t.Errorf("Expected name 'Test Workflow', got %s", workflow.Name)
		}
		
		if workflow.Description != "Description" {
			t.Errorf("Expected description 'Description', got %s", workflow.Description)
		}
		
		if string(workflow.Data) != string(data) {
			t.Errorf("Expected data %s, got %s", string(data), string(workflow.Data))
		}
	})
	
	t.Run("Load non-existent workflow", func(t *testing.T) {
		_, err := store.Load("non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent workflow")
		}
	})
	
	t.Run("Load with empty ID", func(t *testing.T) {
		_, err := store.Load("")
		if err == nil {
			t.Error("Expected error for empty ID")
		}
	})
}

func TestInMemoryStore_Update(t *testing.T) {
	store := NewInMemoryStore()
	
	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	id, err := store.Save("Original Name", "Original Description", data)
	if err != nil {
		t.Fatalf("Failed to save workflow: %v", err)
	}
	
	t.Run("Update existing workflow", func(t *testing.T) {
		newData := json.RawMessage(`{"nodes": [{"id": "1"}], "edges": []}`)
		err := store.Update(id, "Updated Name", "Updated Description", newData)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		
		workflow, err := store.Load(id)
		if err != nil {
			t.Fatalf("Failed to load workflow: %v", err)
		}
		
		if workflow.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %s", workflow.Name)
		}
		
		if workflow.Description != "Updated Description" {
			t.Errorf("Expected description 'Updated Description', got %s", workflow.Description)
		}
		
		if string(workflow.Data) != string(newData) {
			t.Errorf("Expected updated data")
		}
	})
	
	t.Run("Update non-existent workflow", func(t *testing.T) {
		err := store.Update("non-existent", "Name", "Desc", data)
		if err == nil {
			t.Error("Expected error for non-existent workflow")
		}
	})
	
	t.Run("Update with empty ID", func(t *testing.T) {
		err := store.Update("", "Name", "Desc", data)
		if err == nil {
			t.Error("Expected error for empty ID")
		}
	})
	
	t.Run("Update with empty name", func(t *testing.T) {
		err := store.Update(id, "", "Desc", data)
		if err == nil {
			t.Error("Expected error for empty name")
		}
	})
}

func TestInMemoryStore_Delete(t *testing.T) {
	store := NewInMemoryStore()
	
	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	id, err := store.Save("Test Workflow", "Description", data)
	if err != nil {
		t.Fatalf("Failed to save workflow: %v", err)
	}
	
	t.Run("Delete existing workflow", func(t *testing.T) {
		err := store.Delete(id)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		
		// Verify it's deleted
		_, err = store.Load(id)
		if err == nil {
			t.Error("Expected error when loading deleted workflow")
		}
	})
	
	t.Run("Delete non-existent workflow", func(t *testing.T) {
		err := store.Delete("non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent workflow")
		}
	})
	
	t.Run("Delete with empty ID", func(t *testing.T) {
		err := store.Delete("")
		if err == nil {
			t.Error("Expected error for empty ID")
		}
	})
}

func TestInMemoryStore_List(t *testing.T) {
	store := NewInMemoryStore()
	
	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	
	t.Run("Empty store", func(t *testing.T) {
		summaries := store.List()
		if len(summaries) != 0 {
			t.Errorf("Expected empty list, got %d items", len(summaries))
		}
	})
	
	t.Run("Store with workflows", func(t *testing.T) {
		// Save multiple workflows
		id1, _ := store.Save("Workflow 1", "Description 1", data)
		id2, _ := store.Save("Workflow 2", "Description 2", data)
		id3, _ := store.Save("Workflow 3", "Description 3", data)
		
		summaries := store.List()
		
		if len(summaries) != 3 {
			t.Errorf("Expected 3 workflows, got %d", len(summaries))
		}
		
		// Verify all IDs are present
		ids := make(map[string]bool)
		for _, summary := range summaries {
			ids[summary.ID] = true
		}
		
		if !ids[id1] || !ids[id2] || !ids[id3] {
			t.Error("Not all workflow IDs found in list")
		}
	})
}

func TestInMemoryStore_Exists(t *testing.T) {
	store := NewInMemoryStore()
	
	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	id, err := store.Save("Test Workflow", "Description", data)
	if err != nil {
		t.Fatalf("Failed to save workflow: %v", err)
	}
	
	t.Run("Existing workflow", func(t *testing.T) {
		if !store.Exists(id) {
			t.Error("Expected workflow to exist")
		}
	})
	
	t.Run("Non-existent workflow", func(t *testing.T) {
		if store.Exists("non-existent-id") {
			t.Error("Expected workflow to not exist")
		}
	})
}

func TestInMemoryStore_Concurrency(t *testing.T) {
	store := NewInMemoryStore()
	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	
	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			_, err := store.Save("Workflow", "Description", data)
			if err != nil {
				t.Errorf("Failed to save workflow: %v", err)
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	summaries := store.List()
	if len(summaries) != 10 {
		t.Errorf("Expected 10 workflows, got %d", len(summaries))
	}
}
