package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

func main() {
	todos := readTodos()
	for {
		fmt.Println("1. Add todo")
		fmt.Println("2. List todos")
		fmt.Println("3. Delete todo")
		fmt.Println("4. Update todo")
		fmt.Println("5. Exit")
		fmt.Print("Enter your choice: ")
		var choice int
		fmt.Scanln(&choice)
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
			saveTodos(todos)
			os.Exit(0)
		}
	}
}

func readTodos() []Todo {
	file, err := os.Open("todos.json")
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

func saveTodos(todos []Todo) {
	data, err := json.Marshal(todos)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile("todos.json", data, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addTodo() {
	var task string
	fmt.Print("Enter task: ")
	fmt.Scanln(&task)
	todos := readTodos()
	todo := Todo{ID: len(todos) + 1, Task: task}
	todos = append(todos, todo)
	saveTodos(todos)
}

func listTodos(todos []Todo) {
	for _, todo := range todos {
		fmt.Printf("%d. %s\n", todo.ID, todo.Task)
	}
}

func deleteTodo() {
	todos := readTodos()
	listTodos(todos)
	var id int
	fmt.Print("Enter id to delete: ")
	fmt.Scanln(&id)
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTodos(todos)
			break
		}
	}
}

func updateTodo() {
	todos := readTodos()
	listTodos(todos)
	var id int
	fmt.Print("Enter id to update: ")
	fmt.Scanln(&id)
	var task string
	fmt.Print("Enter new task: ")
	fmt.Scanln(&task)
	for i, todo := range todos {
		if todo.ID == id {
			todo.Task = task
			todos[i] = todo
			saveTodos(todos)
			break
		}
	}
}
