package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestReadTodos(t *testing.T) {

	const fileName = "todos.json"

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", fileName)
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
	if err := os.Rename(tmpfile.Name(), fileName); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName) // clean up

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
func TestSaveTodos(t *testing.T) {

	const fileName = "todos.json"

	// Create some test data
	todos := []Todo{
		{ID: 1, Task: "Test task 1"},
		{ID: 2, Task: "Test task 2"},
	}

	// Call the function
	saveTodos(todos)

	// Read the file back
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName) // clean up

	// Read the JSON data and convert it back into a Go object
	var readTodos []Todo
	err = json.Unmarshal(data, &readTodos)
	if err != nil {
		t.Fatal(err)
	}

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
func TestListTodos(t *testing.T) {
	// Capture the output of the function
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create some test data
	todos := []Todo{
		{ID: 1, Task: "Test task 1"},
		{ID: 2, Task: "Test task 2"},
	}

	// Call the function
	listTodos(todos)

	// Restore the original stdout
	w.Close()
	os.Stdout = old

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Expected output
	expectedOutput := "1. Test task 1\n2. Test task 2\n"

	// Check the results
	if buf.String() != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, buf.String())
	}
}
func TestAddTodo(t *testing.T) {

	const fileName = "todos.json"

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Close the file
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Rename the temporary file to "todos.json" to match the function's expectation
	if err := os.Rename(tmpfile.Name(), fileName); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("todos.json") // clean up

	// Simulate user input
	input := "New task\n"
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r

	// Write the input to the pipe
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	// Call the function
	addTodo()

	// Restore the original stdin
	os.Stdin = oldStdin

	// Read the file back
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	// Read the JSON data and convert it back into a Go object
	var readTodos []Todo
	err = json.Unmarshal(data, &readTodos)
	if err != nil {
		t.Fatal(err)
	}

	// Check the results
	if len(readTodos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(readTodos))
	}

	// Check the added todo
	expectedTask := "New task"
	if readTodos[0].Task != expectedTask {
		t.Errorf("Expected task %q, got %q", expectedTask, readTodos[0].Task)
	}
}
func TestDeleteTodo(t *testing.T) {

	const fileName = "todos.json"

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some test data to the temporary file
	todos := []Todo{
		{ID: 1, Task: "Test task 1"},
		{ID: 2, Task: "Test task 2"},
	}

	data, err := json.Marshal(todos)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Rename the temporary file to "todos.json" to match the function's expectation
	if err := os.Rename(tmpfile.Name(), fileName); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName) // clean up

	// Simulate user input
	input := "1\n"
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r

	// Write the input to the pipe
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	// Capture the output of the function
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	// Call the function
	deleteTodo()

	// Restore the original stdin and stdout
	os.Stdin = oldStdin
	wOut.Close()
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, rOut)

	// Read the file back
	data, err = os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	// Read the JSON data and convert it back into a Go object
	var readTodos []Todo
	err = json.Unmarshal(data, &readTodos)
	if err != nil {
		t.Fatal(err)
	}

	// Check the results
	if len(readTodos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(readTodos))
	}

	// Check the remaining todo
	expectedTask := "Test task 2"
	if readTodos[0].Task != expectedTask {
		t.Errorf("Expected task %q, got %q", expectedTask, readTodos[0].Task)
	}
}
func TestUpdateTodo(t *testing.T) {

	const fileName = "todos.json"

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some test data to the temporary file
	todos := []Todo{
		{ID: 1, Task: "Test task 1"},
		{ID: 2, Task: "Test task 2"},
	}

	data, err := json.Marshal(todos)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Rename the temporary file to "todos.json" to match the function's expectation
	if err := os.Rename(tmpfile.Name(), fileName); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fileName) // clean up

	// Simulate user input
	input := "1\nUpdated task\n"
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r

	// Write the input to the pipe
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	// Capture the output of the function
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	// Call the function
	updateTodo()

	// Restore the original stdin and stdout
	os.Stdin = oldStdin
	wOut.Close()
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, rOut)

	// Read the file back
	data, err = os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	// Read the JSON data and convert it back into a Go object
	var readTodos []Todo
	err = json.Unmarshal(data, &readTodos)
	if err != nil {
		t.Fatal(err)
	}

	// Check the results
	if len(readTodos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(readTodos))
	}

	// Check the updated todo
	expectedTask := "Updated task"
	if readTodos[0].Task != expectedTask {
		t.Errorf("Expected task %q, got %q", expectedTask, readTodos[0].Task)
	}
}
func TestGetChoice(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1\n", 1},
		{"2\n", 2},
		{"3\n", 3},
		{"4\n", 4},
		{"5\n", 5},
	}

	for _, test := range tests {
		// Simulate user input
		r, w, _ := os.Pipe()
		oldStdin := os.Stdin
		os.Stdin = r

		// Write the input to the pipe
		go func() {
			w.Write([]byte(test.input))
			w.Close()
		}()

		// Call the function
		choice := getChoice()

		// Restore the original stdin
		os.Stdin = oldStdin

		// Check the result
		if choice != test.expected {
			t.Errorf("Expected choice %d, got %d", test.expected, choice)
		}
	}
}
func TestHandleChoice(t *testing.T) {
	tests := []struct {
		choice   int
		expected string
	}{
		{1, "Enter task: "},
		{2, "1. Test task 1\n2. Test task 2\n"},
		{3, "Enter id to delete: "},
		{4, "Enter id to update: "},
	}

	for _, test := range tests {
		// Create some test data
		todos := []Todo{
			{ID: 1, Task: "Test task 1"},
			{ID: 2, Task: "Test task 2"},
		}

		// Capture the output of the function
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Simulate user input for addTodo, deleteTodo, and updateTodo
		if test.choice == 1 {
			input := "New task\n"
			rIn, wIn, _ := os.Pipe()
			oldStdin := os.Stdin
			os.Stdin = rIn

			go func() {
				wIn.Write([]byte(input))
				wIn.Close()
			}()

			defer func() { os.Stdin = oldStdin }()
		} else if test.choice == 3 {
			input := "1\n"
			rIn, wIn, _ := os.Pipe()
			oldStdin := os.Stdin
			os.Stdin = rIn

			go func() {
				wIn.Write([]byte(input))
				wIn.Close()
			}()

			defer func() { os.Stdin = oldStdin }()
		} else if test.choice == 4 {
			input := "1\nUpdated task\n"
			rIn, wIn, _ := os.Pipe()
			oldStdin := os.Stdin
			os.Stdin = rIn

			go func() {
				wIn.Write([]byte(input))
				wIn.Close()
			}()

			defer func() { os.Stdin = oldStdin }()
		}

		// Call the function
		handleChoice(test.choice, todos)

		// Restore the original stdout
		w.Close()
		os.Stdout = oldStdout

		// Read the captured output
		var buf bytes.Buffer
		io.Copy(&buf, r)

		// Check the results
		if !bytes.Contains(buf.Bytes(), []byte(test.expected)) {
			t.Errorf("Expected output to contain %q, got %q", test.expected, buf.String())
		}
	}
}
func TestPrintMenu(t *testing.T) {
	// Capture the output of the function
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	printMenu()

	// Restore the original stdout
	w.Close()
	os.Stdout = old

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Expected output
	expectedOutput := "1. Add todo\n2. List todos\n3. Delete todo\n4. Update todo\n5. Exit\nEnter your choice: "

	// Check the results
	if buf.String() != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, buf.String())
	}
}
