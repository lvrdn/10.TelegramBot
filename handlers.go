package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Storage interface {
	AddUser(string, int64) error
	AddTask(string, string) (int, error)
	GetAllTasks() ([]*Task, error)
	AddAsigner(string, int) error
	GetTaskWithID(int) (*Task, error)
	SetDoneToTask(id int) error
	CheckUser(string) (*User, error)
}

type Task struct {
	ID       int
	Name     string
	Assigned string
	Owner    string
	Done     bool
}

type User struct {
	ChatID   int64
	UserName string
}

type CommandHandler struct {
	Command  string
	ChatID   int64
	UserName string
	Arg      string
	Storage  Storage
}

func (ch *CommandHandler) AddUser() error {
	err := ch.Storage.AddUser(ch.UserName, ch.ChatID)
	if err != nil {
		return err
	}
	return nil
}

func (ch *CommandHandler) ShowTasks(key string) (string, error) {
	tasks, err := ch.Storage.GetAllTasks()
	var response string

	if err != nil {
		return "", err
	}

	for _, task := range tasks {
		if task.Done {
			continue
		}
		if key == "my" && task.Assigned != ch.UserName {
			continue
		}
		if key == "owner" && task.Owner != ch.UserName {
			continue
		}

		response += fmt.Sprintf("%v. %s by @%s\n", task.ID, task.Name, task.Owner)
		switch task.Assigned {
		case "":
			response += fmt.Sprintf("/assign_%v", task.ID)
		case ch.UserName:
			if key == "my" {
				response += fmt.Sprintf("/unassign_%v /resolve_%v", task.ID, task.ID)
			} else {
				response += fmt.Sprintf("assignee: я\n/unassign_%v /resolve_%v", task.ID, task.ID)
			}

		default:
			response += fmt.Sprintf("assignee: @%s", task.Assigned)
		}

		response += "\n\n"

	}
	if response == "" {
		response = "Нет задач"
	} else {
		response = strings.TrimSuffix(response, "\n\n")
	}

	return response, nil
}

func (ch *CommandHandler) CreateTask() (string, error) {
	id, err := ch.Storage.AddTask(ch.UserName, ch.Arg)
	if err != nil {
		return "", err
	}
	response := fmt.Sprintf(`Задача "%s" создана, id=%v`, ch.Arg, id)
	return response, nil
}

func (ch *CommandHandler) AssignTask() (map[int64]string, error) {

	response := make(map[int64]string)
	taskID, err := strconv.Atoi(ch.Arg)
	if err != nil {
		return nil, fmt.Errorf("task id in command arg is not a number")
	}
	task, err := ch.Storage.GetTaskWithID(taskID)
	if err != nil {
		return nil, err
	}

	if task.Assigned == "" {
		err = ch.Storage.AddAsigner(ch.UserName, taskID)
		if err != nil {
			return nil, err
		}

		owner, err := ch.Storage.CheckUser(task.Owner)
		if err != nil {
			return nil, err
		}

		response[owner.ChatID] = fmt.Sprintf(`Задача "%s" назначена на @%s`, task.Name, task.Assigned)
	} else {
		previousAsigner, err := ch.Storage.CheckUser(task.Assigned)
		if err != nil {
			return nil, err
		}

		err = ch.Storage.AddAsigner(ch.UserName, taskID)
		if err != nil {
			return nil, err
		}

		response[previousAsigner.ChatID] = fmt.Sprintf(`Задача "%s" назначена на @%s`, task.Name, ch.UserName)
	}
	response[ch.ChatID] = fmt.Sprintf(`Задача "%s" назначена на вас`, task.Name)

	return response, nil
}

func (ch *CommandHandler) UnassignTask() (interface{}, error) {
	taskID, err := strconv.Atoi(ch.Arg)
	if err != nil {
		return nil, fmt.Errorf("task id in command arg is not a number")
	}
	task, err := ch.Storage.GetTaskWithID(taskID)
	if err != nil {
		return nil, err
	}

	if task.Assigned == ch.UserName {
		err := ch.Storage.AddAsigner("", task.ID)
		if err != nil {
			return nil, err
		}
		owner, err := ch.Storage.CheckUser(task.Owner)
		if err != nil {
			return nil, err
		}

		response := make(map[int64]string)
		response[ch.ChatID] = "Принято"
		response[owner.ChatID] = fmt.Sprintf(`Задача "%s" осталась без исполнителя`, task.Name)
		return response, nil

	} else {
		return "Задача не на вас", nil
	}
}

func (ch *CommandHandler) ResolveTask() (interface{}, error) {
	taskID, err := strconv.Atoi(ch.Arg)
	if err != nil {
		return nil, fmt.Errorf("task id in command arg is not a number")
	}

	err = ch.Storage.SetDoneToTask(taskID)
	if err != nil {
		return nil, err
	}

	task, err := ch.Storage.GetTaskWithID(taskID)
	if err != nil {
		return nil, err
	}

	owner, err := ch.Storage.CheckUser(task.Owner)
	if err != nil {
		return nil, err
	}

	response := make(map[int64]string)
	response[owner.ChatID] = fmt.Sprintf(`Задача "%s" выполнена @%s`, task.Name, task.Assigned)
	response[ch.ChatID] = fmt.Sprintf(`Задача "%s" выполнена`, task.Name)
	return response, nil
}
