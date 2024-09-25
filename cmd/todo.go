package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"os"
	"time"
)

type Item struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	IsDone      bool      `json:"isDone"`
}

const (
	fileName         = "/Users/naveen.joy/.tasks.json"
	NotPickedStatus  = "Not Picked"
	InProgressStatus = "In Progress"
	DoneStatus       = "Done"
)

type TaskArray []Item

func (taskArray *TaskArray) AddNewItem(name string) {
	newItem := Item{
		Name:   name,
		IsDone: false,
		Status: NotPickedStatus,
	}
	*taskArray = append(*taskArray, newItem)
	taskArray.StoreTasks()
}

func (taskArray *TaskArray) LoadTasks() error {
	data, err := os.ReadFile(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	if len(data) == 0 {
		fmt.Println("File Empty", err)
		return nil
	}
	// Loads data from the file and saves it in taskArray
	err = json.Unmarshal(data, taskArray)
	if err != nil {
		fmt.Println("Error Parsing json file", err)
		return err
	}
	return nil
}

func (taskArray *TaskArray) StoreTasks() {
	taskArray.ListTasks()
	taskJsonBytes, err := json.Marshal(taskArray)

	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	check(err)
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(taskJsonBytes)
	check(err)
}

func (taskArray *TaskArray) DeleteItem(index int) error {
	ls := *taskArray
	if index <= 0 || index > len(ls) {
		return fmt.Errorf("Invalid Index")
	}
	*taskArray = append(ls[:index-1], ls[index:]...)
	taskArray.StoreTasks()
	return nil
}

func (taskArray *TaskArray) CompleteItem(index int) error {
	ls := *taskArray
	if index <= 0 || index > len(ls) {
		return fmt.Errorf("Invalid Index")
	}
	ls[index-1].IsDone = true
	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Status = DoneStatus
	taskArray.StoreTasks()
	return nil
}

func (taskArray *TaskArray) CountPendingTasks() uint16 {
	count := uint16(0)
	for _, task := range *taskArray {
		if !task.IsDone {
			count++
		}
	}
	return count
}

func (taskArray *TaskArray) StartTask(index int) error {
	ls := *taskArray
	if index <= 0 || index > len(ls) {
		return fmt.Errorf("Invalid Index")
	}
	ls[index-1].StartedAt = time.Now()
	ls[index-1].Status = InProgressStatus
	taskArray.StoreTasks()
	return nil
}

func (taskArray *TaskArray) UndoCompletedItem(index int) error {
	ls := *taskArray
	if index <= 0 || index > len(ls) {
		return fmt.Errorf("Invalid Index")
	}
	ls[index-1].Status = InProgressStatus
	ls[index-1].IsDone = false
	taskArray.StoreTasks()
	return nil
}

func (taskArray *TaskArray) ListTasks() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task Name"},
			{Align: simpletable.AlignCenter, Text: "Status"},
			{Align: simpletable.AlignCenter, Text: "Started At"},
			{Align: simpletable.AlignCenter, Text: "Completed At"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
		},
	}
	var cells [][]*simpletable.Cell

	for idx, item := range *taskArray {
		idx++
		task := blue(item.Name)
		done := blue("no")
		var status string
		switch item.Status {
		case NotPickedStatus:
			status = gray(item.Status)
		case InProgressStatus:
			status = amber(item.Status)
		case DoneStatus:
			status = green(item.Status)
		}
		startedAt := formatDate(item.StartedAt)
		completedAt := formatDate(item.CompletedAt)

		if item.IsDone {
			task = green(fmt.Sprintf("\u2705 %s", item.Name))
			done = green("yes")
		}
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: task},
			{Text: status},
			{Text: startedAt},
			{Text: completedAt},
			{Text: done},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}
	var pendingText string
	pendingCount := taskArray.CountPendingTasks()
	if len(*taskArray) == 0 {
		pendingText = green(fmt.Sprintf("Add Tasks to Get Started!"))
	} else if pendingCount > 0 {
		pendingText = red(fmt.Sprintf("You have %d pending todos", taskArray.CountPendingTasks()))
	} else {
		pendingText = green(fmt.Sprintf("All Tasks Are Completed!"))
	}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 6, Text: pendingText},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func formatDate(timeToFormat time.Time) string {
	if timeToFormat.IsZero() {
		return "-"
	}
	return timeToFormat.Format(time.RFC822)
}

var ToDos TaskArray

func check(e error) {
	if e != nil {
		panic(e)
	}
}
