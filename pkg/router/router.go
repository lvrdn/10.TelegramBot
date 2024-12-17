package router

import "taskbot/pkg/handler"

type Router struct {
	handler.Storage
}

func NewRouter(storage handler.Storage) *Router {
	return &Router{
		storage,
	}
}

type CommandManager interface {
	AddUser() error
	CheckUser() error
	ShowTasks(key string) (string, error)
	CreateTask() (string, error)
	AssignTask() (interface{}, error)
	UnassignTask() (interface{}, error)
	ResolveTask() (interface{}, error)
}

func NewCommandManager(command, userName, arg string, chatID int64, storage handler.Storage) CommandManager {
	return handler.NewCommandHandler(command, userName, arg, chatID, storage)
}

func (r *Router) ManageCommand(command, arg, user string, chatID int64) (interface{}, error) {

	cm := NewCommandManager(command, user, arg, chatID, r.Storage)

	err := cm.CheckUser()
	if err != nil {
		return nil, err
	}

	switch command {
	case "start":
		err := cm.AddUser()
		if err != nil {
			return nil, err
		}
		return nil, nil

	case "tasks":
		response, err := cm.ShowTasks("")
		if err != nil {
			return nil, err
		}
		return response, nil

	case "new":
		response, err := cm.CreateTask()
		if err != nil {
			return nil, err
		}
		return response, nil

	case "assign":
		response, err := cm.AssignTask()
		if err != nil {
			return nil, err
		}
		return response, nil

	case "unassign":
		response, err := cm.UnassignTask()
		if err != nil {
			return nil, err
		}
		return response, nil

	case "resolve":
		response, err := cm.ResolveTask()
		if err != nil {
			return nil, err
		}
		return response, nil

	case "my":
		response, err := cm.ShowTasks("my")
		if err != nil {
			return nil, err
		}
		return response, nil

	case "owner":
		response, err := cm.ShowTasks("owner")
		if err != nil {
			return nil, err
		}
		return response, nil

	}
	return "неизвестная команда, выберете команду из списка", nil
}
