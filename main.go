package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Task represents a single to-do item
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskList manages the collection of tasks
type TaskList struct {
	Tasks    []Task `json:"tasks"`
	NextID   int    `json:"next_id"`
	todoFile string
}

// NewTaskList creates a new TaskList
func NewTaskList(filename string) *TaskList {
	return &TaskList{
		Tasks:    []Task{},
		NextID:   1,
		todoFile: filename,
	}
}

// Add new task to tasklist
func (tl *TaskList) AddTask(title string) Task {
	task := Task{
		ID:        tl.NextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	tl.Tasks = append(tl.Tasks, task)
	tl.NextID++
	return task
}

// Delete function by task id

func (tl *TaskList) DeleteTask(id int) error {
	for i, task := range tl.Tasks {
		if task.ID == id {
			// remove the task at index id
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

// ToggleTaskStatus changes the state of the task
func (tl *TaskList) ToggleTaskStatus(id int) error {
	for i, task := range tl.Tasks {
		if task.ID == id {
			tl.Tasks[i].Completed = !tl.Tasks[i].Completed
			return nil
		}
	}
	return errors.New("task not found")
}

// GetTask get task from id
func (tl *TaskList) GetTask(id int) (Task, error) {
	for _, task := range tl.Tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return Task{}, errors.New("task not found")
}

// SaveToFile saves task to json file
func (tl *TaskList) SaveToFile() error {
	data, err := json.MarshalIndent(tl, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(tl.todoFile, data, 0o644)
}

// LoadFromFile -> loads tasks from a json file
func (tl *TaskList) LoadFromFile() error {
	data, err := ioutil.ReadFile(tl.todoFile)
	if err != nil {
		// if file doesn't exist return without error
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, tl)
}

// PrintTasks -> print tasks in a formatted way
func (tl *TaskList) PrintTasks() {
	if len(tl.Tasks) == 0 {
		fmt.Println("No Tasks found")
		return
	}

	fmt.Println("ID | Complete | Title")
	fmt.Println("...................")
	for _, task := range tl.Tasks {
		status := " "
		if task.Completed {
			status = "✓"
		}
		fmt.Printf("%2d | [%s] | %s\n", task.ID, status, task.Title)
	}
}

func displayHelp() {
	fmt.Println("Todo List Application")
	fmt.Println("Commands:")
	fmt.Println("  list              - List all tasks")
	fmt.Println("  add <title>       - Add a new task")
	fmt.Println("  complete <id>     - Toggle task completion status")
	fmt.Println("  delete <id>       - Delete a task")
	fmt.Println("  help              - Display this help message")
	fmt.Println("  exit              - Exit the application")
}

func main() {
	todoList := NewTaskList("todos.json")
	err := todoList.LoadFromFile()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}
	fmt.Println("Todo List Application")
	fmt.Println("Type 'help' to see available commands")
	for {
		fmt.Print("> ")
		var input string
		// Use Scanln to read a line of input
		input, _ = readLine()

		// Skip empty input
		if input == "" {
			continue
		}

		// Split the input to get command and arguments
		parts := strings.SplitN(input, " ", 2)
		command := parts[0]

		var args string
		if len(parts) > 1 {
			args = parts[1]
		}

		switch command {
		case "list":
			todoList.PrintTasks()
		case "add":
			if args == "" {
				fmt.Println("Error: Task title required")
				continue
			}
			task := todoList.AddTask(args)
			fmt.Printf("Added task: %s (ID: %d)\n", task.Title, task.ID)
			todoList.SaveToFile()
		case "complete":
			if args == "" {
				fmt.Println("Error: Task ID required")
				continue
			}
			id, err := strconv.Atoi(args)
			if err != nil {
				fmt.Println("Error: Invalid task ID")
				continue
			}
			err = todoList.ToggleTaskStatus(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Toggled status of task %d\n", id)
			todoList.SaveToFile()
		case "delete":
			if args == "" {
				fmt.Println("Error: Task ID required")
				continue
			}
			id, err := strconv.Atoi(args)
			if err != nil {
				fmt.Println("Error: Invalid task ID")
				continue
			}
			err = todoList.DeleteTask(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Deleted task %d\n", id)
			todoList.SaveToFile()
		case "help":
			displayHelp()
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

// Helper function to read a full line of input
func readLine() (string, error) {
	var input string
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// Trim spaces and newline characters
	return strings.TrimSpace(input), nil
}
