package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID                                   int
	Title, Description, Priority, Status string
}

func (task Task) show() {
	fmt.Printf("#%d | %s | %s | %s | %s\n", task.ID, task.Title, task.Priority, task.Status, task.Description)
}

// Интерфейс перечисляет действия с задачами.
type TaskStorage interface {
	Add(title, description, priority string) error
	All() []Task
	Find(id int) (*Task, error)
	Delete(id int) error
}
type MemoryStorage struct {
	tasks  []Task
	nextID int
}

func validPriority(product string) bool {
	return product == "низкий" || product == "средний" || product == "высокий"
}

func validStatus(storage string) bool {
	return storage == "новая" || storage == "в работе" || storage == "завершена"
}

func validate(title, priority string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("название не может быть пустым")
	}
	if !validPriority(priority) {
		return errors.New("приоритет: низкий, средний или высокий")
	}
	return nil
}

func (storage *MemoryStorage) Add(title, description, priority string) error {
	if inputError := validate(title, priority); inputError != nil {
		return inputError
	}
	storage.nextID = storage.nextID + 1
	storage.tasks = append(storage.tasks, Task{storage.nextID, title, description, priority, "новая"})
	return nil
}

func (storage *MemoryStorage) All() []Task {
	return storage.tasks
}

// Ищем задачу по номеру и возвращаем указатель на неё.
func (storage *MemoryStorage) Find(id int) (*Task, error) {
	for index := range storage.tasks {
		if storage.tasks[index].ID == id {
			return &storage.tasks[index], nil
		}
	}
	return nil, errors.New("задача не найдена")
}

func (storage *MemoryStorage) Delete(id int) error {
	for index, task := range storage.tasks {
		if task.ID == id {
			storage.tasks = append(storage.tasks[:index], storage.tasks[index+1:]...)
			return nil
		}
	}
	return errors.New("задача не найдена")
}

var input = bufio.NewScanner(os.Stdin)

func read(prompt string) string {
	fmt.Print(prompt)
	if !input.Scan() {
		return "0"
	}
	return strings.TrimSpace(input.Text())
}

func getTask(storage TaskStorage) (*Task, error) {
	id, inputError := strconv.Atoi(read("ID задачи: "))
	if inputError != nil {
		return nil, errors.New("неверный ID")
	}
	return storage.Find(id)
}

func main() {
	var storage TaskStorage = &MemoryStorage{}
	for {
		fmt.Println()
		fmt.Println("1. Создать задачу\n2. Показать все задачи\n3. Изменить статус\n4. Редактировать задачу\n5. Удалить задачу\n6. Поиск и фильтрация\n7. Статистика\n0. Выход")
		switch read("Выберите действие: ") {
		case "1":
			title := read("Название: ")
			description := read("Описание: ")
			priority := strings.ToLower(read("Приоритет (низкий/средний/высокий): "))
			if inputError := storage.Add(title, description, priority); inputError != nil {
				fmt.Println(inputError)
			} else {
				fmt.Println("Задача создана")
			}
		case "2":
			if len(storage.All()) == 0 {
				fmt.Println("Задач нет")
			}
			for _, task := range storage.All() {
				task.show()
			}
		case "3":
			task, inputError := getTask(storage)
			if inputError != nil {
				fmt.Println(inputError)
				continue
			}
			status := strings.ToLower(read("Статус (новая/в работе/завершена): "))
			if !validStatus(status) {
				fmt.Println("Неверный статус")
				continue
			}
			task.Status = status
			fmt.Println("Статус изменен")
		case "4":
			task, inputError := getTask(storage)
			if inputError != nil {
				fmt.Println(inputError)
				continue
			}
			title := read("Название: ")
			description := read("Описание: ")
			priority := strings.ToLower(read("Приоритет: "))
			if inputError := validate(title, priority); inputError != nil {
				fmt.Println(inputError)
				continue
			}
			task.Title = title
			task.Description = description
			task.Priority = priority
			fmt.Println("Задача изменена")
		case "5":
			id, inputError := strconv.Atoi(read("ID задачи: "))
			if inputError != nil {
				fmt.Println("Неверный ID")
				continue
			}
			if inputError := storage.Delete(id); inputError != nil {
				fmt.Println(inputError)
			} else {
				fmt.Println("Задача удалена")
			}
		case "6":
			mode := read("Поиск: 1 — название, 2 — приоритет, 3 — статус: ")
			value := strings.ToLower(read("Значение: "))
			found := false
			for _, task := range storage.All() {
				match := mode == "1" && strings.Contains(strings.ToLower(task.Title), value) || mode == "2" && strings.EqualFold(task.Priority, value) || mode == "3" && strings.EqualFold(task.Status, value)
				if match {
					task.show()
					found = true
				}
			}
			if !found {
				fmt.Println("Задачи не найдены")
			}
		case "7":
			counts := map[string]int{}
			tasks := storage.All()
			for _, task := range tasks {
				counts[task.Status] = counts[task.Status] + 1
			}
			percent := 0.0
			if len(tasks) > 0 {
				percent = float64(counts["завершена"]) / float64(len(tasks)) * 100
			}
			fmt.Printf("Всего: %d\nНовых: %d\nВ работе: %d\nЗавершенных: %d\nВыполнение: %.0f%%\n", len(tasks), counts["новая"], counts["в работе"], counts["завершена"], percent)
		case "0":
			return
		default:
			fmt.Println("Неизвестный пункт")
		}
	}
}
