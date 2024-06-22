package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsDone      bool      `json:"isDone"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Todos []Todo

var todos Todos

func addTodo(title, description string) {
	loadTodosFromFile()
	var id int

	if len(todos) == 0 {
		id = 1
	} else {
		id = todos[len(todos)-1].ID + 1
	}
	todo := Todo{
		ID:          id,
		Title:       title,
		Description: description,
		IsDone:      false,
		CreatedAt:   time.Now(),
	}

	todos = append(todos, todo)
	saveTodosToFile()
	fmt.Printf("New Todo added:\nID: %d\nTitle: %s\nDescription: %s\nCreated At: %s\n", todo.ID, todo.Title, todo.Description, todo.CreatedAt.Format(time.RFC1123))
}

func saveTodosToFile() {
	homedir := getHomeDir()

	todofilePath := filepath.Join(homedir, ".todos.json")

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling todos to JSON: %v", err)
	}

	err = os.WriteFile(todofilePath, data, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}
}

func loadTodosFromFile() {
	homedir := getHomeDir()

	todofilePath := filepath.Join(homedir, ".todos.json")

	data, err := os.ReadFile(todofilePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.WriteFile(todofilePath, []byte("[]"), 0644)
			if err != nil {
				log.Fatalf("Error creating todos file: %v", err)
			}
			todos = make(Todos, 0)
			return
		}
		log.Fatalf("Error reading todos file: %v", err)
	}

	err = json.Unmarshal(data, &todos)
	if err != nil {
		log.Fatalf("Error parsing JSON file: %v", err)
	}
}

func updateTodo(id int, title, description string, isDone bool) {
	loadTodosFromFile()

	found := false
	for idx, todo := range todos {
		if todo.ID == id {
			found = true
			if title != "" {
				todos[idx].Title = title
			}
			if description != "" {
				todos[idx].Description = description
			}
			todos[idx].IsDone = isDone

			fmt.Printf("Updated Todo with ID %d:\n", todo.ID)
			fmt.Printf("Title: %s\n", todos[idx].Title)
			fmt.Printf("Description: %s\n", todos[idx].Description)
			fmt.Printf("IsDone: %v\n", todos[idx].IsDone)

			break
		}
	}

	if !found {
		log.Fatalf("Todo with ID %d not found", id)
	}

	saveTodosToFile()
	printTodosTable()
}

func listTodos() {
	loadTodosFromFile()
	printTodosTable()
}

func deleteTodo(id int) {
	loadTodosFromFile()

	idx := -1
	for i, t := range todos {
		if t.ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		log.Fatalf("Todo with ID %d not found", id)
	}

	todos = append(todos[:idx], todos[idx+1:]...)
	saveTodosToFile()
	fmt.Printf("Todo with ID %d has been deleted.\n", id)
	fmt.Println("Updated list of todos.")
	printTodosTable()
}

func printTodosTable() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)

	fmt.Fprintln(w, "ID\t|\tTitle\t|\tDescription\t|\tCreated At\t|\t\tDone")
	fmt.Fprintln(w, "------------------------------------------------------------")

	for _, todo := range todos {
		doneStatus := color.RedString("false")
		if todo.IsDone {
			doneStatus = color.GreenString("true")
		}

		createdAt := todo.CreatedAt.Format("Jan 02 2006 15:04:05")

		fmt.Fprintf(w, "%d\t%s\t%s\t|\t%s\t|\t\t%s\n", todo.ID, todo.Title, todo.Description, createdAt, doneStatus)
	}

	w.Flush()
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting Home dir: %v", err)
	}

	return homeDir
}

