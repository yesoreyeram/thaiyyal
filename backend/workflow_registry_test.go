package workflow

import (
	"encoding/json"
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	registry := NewWorkflowRegistry()

	data := json.RawMessage(`{"nodes": [], "edges": []}`)

	tests := []struct {
		name         string
		workflowName string
		description  string
		data         json.RawMessage
		wantErr      bool
	}{
		{
			name:         "Valid workflow",
			workflowName: "Test Workflow",
			description:  "A test workflow",
			data:         data,
			wantErr:      false,
		},
		{
			name:         "Empty name",
			workflowName: "",
			description:  "Description",
			data:         data,
			wantErr:      true,
		},
		{
			name:         "Empty data",
			workflowName: "Test",
			description:  "Description",
			data:         json.RawMessage{},
			wantErr:      true,
		},
		{
			name:         "Invalid JSON data",
			workflowName: "Test",
			description:  "Description",
			data:         json.RawMessage(`{invalid json`),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := registry.Register(tt.workflowName, tt.description, tt.data)

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

func TestRegistry_Get(t *testing.T) {
	registry := NewWorkflowRegistry()

	data := json.RawMessage(`{"nodes": [{"id": "1"}], "edges": []}`)
	id, err := registry.Register("Test Workflow", "Description", data)
	if err != nil {
		t.Fatalf("Failed to register workflow: %v", err)
	}

	t.Run("Get existing workflow", func(t *testing.T) {
		workflow, err := registry.Get(id)
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

	t.Run("Get non-existent workflow", func(t *testing.T) {
		_, err := registry.Get("non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent workflow")
		}
	})

	t.Run("Get with empty ID", func(t *testing.T) {
		_, err := registry.Get("")
		if err == nil {
			t.Error("Expected error for empty ID")
		}
	})
}

func TestRegistry_Update(t *testing.T) {
	registry := NewWorkflowRegistry()

	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	id, err := registry.Register("Original Name", "Original Description", data)
	if err != nil {
		t.Fatalf("Failed to register workflow: %v", err)
	}

	t.Run("Update existing workflow", func(t *testing.T) {
		newData := json.RawMessage(`{"nodes": [{"id": "1"}], "edges": []}`)
		err := registry.Update(id, "Updated Name", "Updated Description", newData)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		workflow, err := registry.Get(id)
		if err != nil {
			t.Fatalf("Failed to get workflow: %v", err)
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
		err := registry.Update("non-existent", "Name", "Desc", data)
		if err == nil {
			t.Error("Expected error for non-existent workflow")
		}
	})

	t.Run("Update with empty ID", func(t *testing.T) {
		err := registry.Update("", "Name", "Desc", data)
		if err == nil {
			t.Error("Expected error for empty ID")
		}
	})

	t.Run("Update with empty name", func(t *testing.T) {
		err := registry.Update(id, "", "Desc", data)
		if err == nil {
			t.Error("Expected error for empty name")
		}
	})
}

func TestRegistry_Unregister(t *testing.T) {
	registry := NewWorkflowRegistry()

	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	id, err := registry.Register("Test Workflow", "Description", data)
	if err != nil {
		t.Fatalf("Failed to register workflow: %v", err)
	}

	t.Run("Unregister existing workflow", func(t *testing.T) {
		err := registry.Unregister(id)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		// Verify it's deleted
		_, err = registry.Get(id)
		if err == nil {
			t.Error("Expected error when getting unregistered workflow")
		}
	})

	t.Run("Unregister non-existent workflow", func(t *testing.T) {
		err := registry.Unregister("non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent workflow")
		}
	})

	t.Run("Unregister with empty ID", func(t *testing.T) {
		err := registry.Unregister("")
		if err == nil {
			t.Error("Expected error for empty ID")
		}
	})
}

func TestRegistry_List(t *testing.T) {
	registry := NewWorkflowRegistry()

	data := json.RawMessage(`{"nodes": [], "edges": []}`)

	t.Run("Empty registry", func(t *testing.T) {
		summaries := registry.List()
		if len(summaries) != 0 {
			t.Errorf("Expected empty list, got %d items", len(summaries))
		}
	})

	t.Run("Registry with workflows", func(t *testing.T) {
		// Register multiple workflows
		id1, _ := registry.Register("Workflow 1", "Description 1", data)
		id2, _ := registry.Register("Workflow 2", "Description 2", data)
		id3, _ := registry.Register("Workflow 3", "Description 3", data)

		summaries := registry.List()

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

func TestRegistry_Has(t *testing.T) {
	registry := NewWorkflowRegistry()

	data := json.RawMessage(`{"nodes": [], "edges": []}`)
	id, err := registry.Register("Test Workflow", "Description", data)
	if err != nil {
		t.Fatalf("Failed to register workflow: %v", err)
	}

	t.Run("Existing workflow", func(t *testing.T) {
		if !registry.Has(id) {
			t.Error("Expected workflow to exist")
		}
	})

	t.Run("Non-existent workflow", func(t *testing.T) {
		if registry.Has("non-existent-id") {
			t.Error("Expected workflow to not exist")
		}
	})
}

func TestRegistry_Concurrency(t *testing.T) {
	registry := NewWorkflowRegistry()
	data := json.RawMessage(`{"nodes": [], "edges": []}`)

	// Test concurrent registrations
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			_, err := registry.Register("Workflow", "Description", data)
			if err != nil {
				t.Errorf("Failed to register workflow: %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	summaries := registry.List()
	if len(summaries) != 10 {
		t.Errorf("Expected 10 workflows, got %d", len(summaries))
	}
}
