package main

type Router struct {
	Data Storage
}

func (r *Router) RouteManage(command, arg, user string, chatID int64) (interface{}, error) {

	handler := &CommandHandler{
		Command:  command,
		ChatID:   chatID,
		Arg:      arg,
		UserName: user,
		Storage:  r.Data,
	}

	_, err := handler.Storage.CheckUser(handler.UserName)
	if err != nil {
		err := handler.Storage.AddUser(handler.UserName, handler.ChatID)
		if err != nil {
			return nil, err
		}
	}

	switch command {
	case "start":
		err := handler.AddUser()
		if err != nil {
			return nil, err
		}
		return nil, nil

	case "tasks":
		response, _ := handler.ShowTasks("")
		return response, nil

	case "new":
		response, err := handler.CreateTask()
		if err != nil {
			return "", err
		}
		return response, nil

	case "assign":
		response, err := handler.AssignTask()
		if err != nil {
			return nil, err
		}
		return response, nil

	case "unassign":
		response, err := handler.UnassignTask()
		if err != nil {
			return nil, err
		}
		return response, nil

	case "resolve":
		response, err := handler.ResolveTask()
		if err != nil {
			return nil, err
		}
		return response, nil

	case "my":
		response, err := handler.ShowTasks("my")
		if err != nil {
			return nil, err
		}
		return response, nil

	case "owner":
		response, err := handler.ShowTasks("owner")
		if err != nil {
			return nil, err
		}
		return response, nil

	}
	return "unknown command", nil
}
