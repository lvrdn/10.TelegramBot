package storage

import (
	"fmt"
	"sync"
	"taskbot/pkg/handler"
)

type StorageMemory struct {
	Users   []*handler.User
	Tasks   []*handler.Task
	MuUsers *sync.RWMutex
	MuTasks *sync.RWMutex
}

func NewStorage() (*StorageMemory, error) {
	return &StorageMemory{
		Tasks:   make([]*handler.Task, 0, 10),
		Users:   make([]*handler.User, 0, 10),
		MuUsers: &sync.RWMutex{},
		MuTasks: &sync.RWMutex{},
	}, nil
}

func (st *StorageMemory) AddUser(userName string, chatID int64) error {

	st.MuUsers.Lock()
	defer st.MuUsers.Unlock()
	st.Users = append(st.Users, &handler.User{
		UserName: userName,
		ChatID:   chatID,
	})

	return nil
}

func (st *StorageMemory) CheckUser(userName string) (*handler.User, error) {
	st.MuUsers.RLock()
	defer st.MuUsers.RUnlock()

	for _, user := range st.Users {
		if user.UserName == userName {
			return user, nil
		}
	}
	return nil, st.GetErrorNoUser()
}

func (st *StorageMemory) GetAllTasks() ([]*handler.Task, error) {
	st.MuTasks.RLock()
	defer st.MuTasks.RUnlock()
	tasks := st.Tasks
	return tasks, nil
}

func (st *StorageMemory) AddTask(user, task string) (int, error) {
	st.MuTasks.Lock()
	defer st.MuTasks.Unlock()
	id := len(st.Tasks) + 1

	new := &handler.Task{
		ID:    id,
		Name:  task,
		Owner: user,
	}

	st.Tasks = append(st.Tasks, new)

	return id, nil
}

func (st *StorageMemory) AddAsigner(user string, id int) error {
	st.MuTasks.Lock()
	defer st.MuTasks.Unlock()
	for i := 0; i < len(st.Tasks); i++ {
		if st.Tasks[i].ID == id {
			st.Tasks[i].Assigned = user
			return nil
		}
	}
	return fmt.Errorf("assigned is not updated, no task with this id")
}

func (st *StorageMemory) GetTaskWithID(id int) (*handler.Task, error) {
	st.MuTasks.RLock()
	defer st.MuTasks.RUnlock()
	for i := 0; i < len(st.Tasks); i++ {
		if st.Tasks[i].ID == id {
			return st.Tasks[i], nil
		}
	}
	return nil, fmt.Errorf("no task with this id")
}

func (st *StorageMemory) SetDoneToTask(id int) error {
	st.MuTasks.Lock()
	defer st.MuTasks.Unlock()
	for i := 0; i < len(st.Tasks); i++ {
		if st.Tasks[i].ID == id {
			st.Tasks[i].Done = true
			return nil
		}
	}
	return fmt.Errorf("no task with this id")
}

func (st *StorageMemory) GetErrorNoUser() error {
	return fmt.Errorf("there is no user with this username")
}
