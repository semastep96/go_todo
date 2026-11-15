package main

import (
	"errors"
	"fmt"
)

type Task struct {
	Title    string
	Status   string
	Priority string
}
type TaskList []Task

func main() {
	taskList := TaskList{}
	AddTask(&taskList, Task{
		Title:    "Task1",
		Status:   "todo",
		Priority: "low",
	})
	AddTask(&taskList, Task{
		Title:    "Task2",
		Status:   "progress",
		Priority: "high",
	})
	AddTask(&taskList, Task{
		Title:    "Task3",
		Status:   "done",
		Priority: "high",
	})
	PrintTasks(taskList)
	fmt.Println(AddTask(&taskList, Task{}))
	fmt.Println(AddTask(&taskList, Task{Title: "ok title"}))
	fmt.Println(AddTask(&taskList, Task{Title: "ok title", Status: "todo"}))
	fmt.Println(AddTask(&taskList, taskList[0]))
	ChangeTaskByTitle(taskList, taskList[0].Title, "done", taskList[0].Priority)
	PrintTasks(taskList)
	fmt.Println(ChangeTaskByTitle(taskList, taskList[0].Title, "todo", "random"))
	fmt.Println(ChangeTaskByTitle(taskList, "random", "todo", "random"))
	RemoveTaskByTitle(&taskList, taskList[0].Title)
	PrintTasks(taskList)
	fmt.Println(RemoveTaskByTitle(&taskList, "some task"))
	fmt.Println(RemoveTaskByTitle(&taskList, "Task2"))
}

func AddTask(tl *TaskList, task Task) error {
	err := CheckTaskFields(task)
	if err != nil {
		return err
	}
	if IsUniqueTaskName(*tl, task) {
		*tl = append(*tl, task)
		return nil
	}
	return errors.New("task already exists")
}

func CheckTaskFields(task Task) error {
	if task.Title == "" {
		return errors.New("empty title")
	}

	switch task.Status {
	case "todo", "progress", "done":
	default:
		return errors.New("unknown status")
	}

	switch task.Priority {
	case "low", "high":
	default:
		return errors.New("unknown priority")
	}
	return nil
}

func IsUniqueTaskName(tl TaskList, task Task) bool {
	for _, t := range tl {
		if t.Title == task.Title {
			return false
		}
	}
	return true
}

func RemoveTaskByTitle(tl *TaskList, title string) error {
	for i, t := range *tl {
		isProgressTask := t.Status == "progress"
		isSameTitle := t.Title == title
		if isSameTitle && isProgressTask {
			return errors.New("can't remove task in progress")
		}
		if isSameTitle {
			*tl = append((*tl)[:i], (*tl)[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

func ChangeTaskByTitle(tl TaskList, title string, status string, priority string) error {
	for i, task := range tl {
		if task.Title == title {
			task.Priority = priority
			task.Status = status
			err := CheckTaskFields(task)
			if err != nil {
				return err
			}
			tl[i] = task
			return nil
		}
	}
	return errors.New("task not found")
}

func PrintTasks(tl TaskList) {
	for _, task := range tl {
		fmt.Println(task.Title, "-", task.Status)
	}
}
