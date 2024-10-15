package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

const fileName = "todos.json"

func main() {
	for {
		todos := readTodos()
		printMenu()
		choice := getChoice()
		handleChoice(choice, todos)
	}
}

// getChoice reads an integer input from the user and returns it.
// It uses fmt.Scanln to capture the input and stores it in the variable 'choice'.
// Returns:
//
//	int: The integer value entered by the user.
func getChoice() int {
	var choice int
	fmt.Scanln(&choice)
	return choice
}

// handleChoice takes an integer choice and a slice of todos as input.
func handleChoice(choice int, todos []Todo) {
	switch choice {
	case 1:
		addTodo()
	case 2:
		listTodos(todos)
	case 3:
		deleteTodo()
	case 4:
		updateTodo()
	case 5:
		os.Exit(0)
	}
}

// printMenu prints the menu options to the console.
func printMenu() {
	fmt.Println("1. Add todo")
	fmt.Println("2. List todos")
	fmt.Println("3. Delete todo")
	fmt.Println("4. Update todo")
	fmt.Println("5. Exit")
	fmt.Print("Enter your choice: ")
}

// readTodos reads the todos from the file and returns them.
func readTodos() []Todo {
	file, err := os.Open(fileName)
	if err != nil {
		return []Todo{}
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return []Todo{}
	}

	var todos []Todo

	err = json.Unmarshal(data, &todos)
	if err != nil {
		return []Todo{}
	}
	return todos
}

// saveTodos saves the todos to the file.
func saveTodos(todos []Todo) {
	data, err := json.Marshal(todos)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// addTodo adds a new todo to the list.
func addTodo() {
	var task string
	fmt.Print("Enter task: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	task = scanner.Text()
	fmt.Println("The task is " + task)
	todos := readTodos()
	todo := Todo{ID: len(todos) + 1, Task: task}
	todos = append(todos, todo)
	saveTodos(todos)
}

// listTodos lists all the todos to the console.
func listTodos(todos []Todo) {
	for _, todo := range todos {
		fmt.Printf("%d. %s\n", todo.ID, todo.Task)
	}
}

// deleteTodo deletes a todo from the list.
func deleteTodo() {
	todos := readTodos()
	listTodos(todos)
	fmt.Print("Enter id to delete: ")
	var id int
	fmt.Scanln(&id)
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTodos(todos)
			break
		}
	}
}

// updateTodo updates a todo in the list.
func updateTodo() {
	todos := readTodos()
	listTodos(todos)
	fmt.Print("Enter id to update: ")
	var id int
	fmt.Scanln(&id)
	var task string
	fmt.Print("Enter new task: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	task = scanner.Text()
	for i, todo := range todos {
		if todo.ID == id {
			todo.Task = task
			todos[i] = todo
			saveTodos(todos)
			break
		}
	}
}
