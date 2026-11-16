package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Task struct {
	Title    string
	Status   string
	Priority string
}
type TaskList []Task

func main() {
	taskList := TaskList{}
	fmt.Println("0/ Добро пожаловать в todo 2000!")
	fmt.Println("Выход exit, помощь help")
	loadErr := Load(&taskList, "tasks.gob")
	if loadErr != nil {
		fmt.Println("Ошибка при загрузке задач", loadErr)
		return
	}
	handleIO(&taskList)
	fmt.Println("Спасибо что не ctrl+c :D")
	fmt.Println("До скорой встречи 0/")
}

func Load(tasks *TaskList, filename string) error {
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	var loaded []Task
	decodeErr := gob.NewDecoder(file).Decode(&loaded)
	closeErr := file.Close()
	if decodeErr != nil {
		return decodeErr
	}
	if closeErr != nil {
		return closeErr
	}

	*tasks = loaded
	return nil
}

func Save(tl TaskList, fileName string) error {
	file, createErr := os.Create(fileName)
	if createErr != nil {
		return createErr
	}
	encodeErr := gob.NewEncoder(file).Encode(tl)
	closeErr := file.Close()
	if encodeErr != nil {
		return encodeErr
	}
	return closeErr
}

func HandleDefaultSave(tl TaskList) {
	err := Save(tl, "tasks.gob")
	if err != nil {
		fmt.Println("Не удалось сохранить задачи:", err)
	} else {
		fmt.Println("Задачи сохранены")
	}
}

func handleIO(tl *TaskList) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			fmt.Println(err)
			continue
		}
		parts := ParseUserInput(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]
		switch command {
		case "exit":
			return
		case "help":
			PrintHelp()
		case "list":
			PrintTasks(*tl)
		case "search":
			if len(args) != 1 {
				fmt.Println("Ожидается подстрока для поиска, если подстрока содержит пробел используйте кавычки")
				continue
			}
			SearchTask(*tl, args[0])
		case "add":
			if len(args) < 2 {
				fmt.Println("Необходимо указать название и приоритет")
				continue
			}
			err = AddTask(tl, Task{Title: args[0], Status: "todo", Priority: args[1]})
			if err != nil {
				fmt.Println(err)
				continue
			}
			HandleDefaultSave(*tl)
		case "change":
			if len(args) < 3 {
				fmt.Println("Необходимо указать название статус и приоритет")
				continue
			}
			err := ChangeTaskByTitle(*tl, args[0], args[1], args[2])
			if err != nil {
				fmt.Println(err)
				continue
			}
			HandleDefaultSave(*tl)
		case "rm":
			if len(args) < 1 {
				fmt.Println("Необходимо указать название задачи для удаления")
				continue
			}
			err := RemoveTaskByTitle(tl, args[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
			HandleDefaultSave(*tl)
		default:
			fmt.Println("Неизвестная команда. Введите help для списка команд")
		}
	}
}

func PrintHelp() {
	fmt.Println("Commands:")
	fmt.Println("exit - выйти из программы")
	fmt.Println("help - показать команды")
	fmt.Println("list - показать все задачи")
	fmt.Println("search <substr> - задачи в названии которых есть подстрока")
	fmt.Println("add <title> <priority: [low/high]> - добавляет задачу с конкретным приоритетом")
	fmt.Println("change <title> <status: [todo/progress/done]> <priority: [low/high]> - меняет статус и приоритет задачи с конкретным названием")
	fmt.Println("rm <title> - удаляет задачу с конкретным названием")
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

func SearchTask(tl TaskList, query string) {
	found := false
	fmt.Println("Совпадения:")
	for _, task := range tl {
		if strings.Contains(task.Status, query) {
			found = true
			fmt.Print("  - ")
			PrintTask(task)
		}
	}
	if found {
		return
	}
	fmt.Println("Найдено 0")
}

func PrintTasks(tl TaskList) {
	if len(tl) == 0 {
		fmt.Println("Список задач пуст")
		return
	}
	for i, task := range tl {
		fmt.Print("  ", i+1, " | ")
		PrintTask(task)
	}
}

func PrintTask(t Task) {
	fmt.Println("Заголовок:", t.Title, "Статус:", t.Status, "Приоритет:", t.Priority)
}

func ParseUserInput(userInput string) []string {
	args := []string{}
	word := ""
	inside := false
	for _, r := range userInput {
		switch r {
		case '"', '\'':
			inside = !inside
			continue
		case ' ':
			if inside {
				word += string(r)
				continue
			}
			if word != "" {
				args = append(args, word)
				word = ""
			}
		default:
			word += string(r)
		}
	}
	if word != "" {
		args = append(args, word)
	}
	return args
}
