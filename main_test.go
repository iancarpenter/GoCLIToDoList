package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestReadTodos(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "todos.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some test data to the temporary file
	todos := []Todo{
		{ID: 1, Task: "Test task 1"},
		{ID: 2, Task: "Test task 2"},
	}

	// Returns the JSON-encoded byte array and an error if any occurred during the encoding process.
	data, err := json.Marshal(todos)
	if err != nil {
		t.Fatal(err)
	}

	// Write the JSON-encoded byte array to the temporary file
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}

	// Close the file
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Rename the temporary file to "todos.json" to match the function's expectation
	if err := os.Rename(tmpfile.Name(), "todos.json"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("todos.json") // clean up

	// Call the function
	readTodos := readTodos()

	// Check the results
	if len(readTodos) != len(todos) {
		t.Errorf("Expected %d todos, got %d", len(todos), len(readTodos))
	}

	// Check each todo
	for i, todo := range readTodos {
		if todo.ID != todos[i].ID || todo.Task != todos[i].Task {
			t.Errorf("Expected todo %v, got %v", todos[i], todo)
		}
	}
}
