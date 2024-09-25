package main

import (
	"bufio"
	"flag"
	"fmt"
	todo "github.com/NaveenJoyGit/go-todo-app/cmd"
	"io"
	"os"
	"strings"
)

func main() {
	addItemPtr := flag.Bool("a", false, "Item to add to the to-do list")
	deleteItem := flag.Int("d", 0, "Item to delete from the to-do list")
	complete := flag.Int("c", 0, "Item to complete from the to-do list")
	startItem := flag.Int("s", 0, "Item to start from the to-do list")
	undoCompleteItem := flag.Int("u", 0, "Item to undo from the to-do list")
	flag.Parse()
	todos := &todo.TaskArray{}
	err := todos.LoadTasks()
	if err != nil {
		_ = fmt.Errorf(err.Error())
		os.Exit(1)
	}
	switch {
	case *addItemPtr:
		input, err := readTextFromInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		todos.AddNewItem(input)
	case *deleteItem > 0:
		err := todos.DeleteItem(*deleteItem)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case *complete > 0:
		err := todos.CompleteItem(*complete)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case *startItem > 0:
		err := todos.StartTask(*startItem)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case *undoCompleteItem > 0:
		err := todos.UndoCompletedItem(*undoCompleteItem)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		todos.ListTasks()
	}
	//todos.LoadTasks()
}

func readTextFromInput(r io.Reader, arg ...string) (string, error) {
	if len(arg) > 1 {
		strings.Join(arg, " ")
	}
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	text := scanner.Text()
	if len(text) == 0 {
		return "", fmt.Errorf("empty todo not allowed")
	}
	return text, nil

}
