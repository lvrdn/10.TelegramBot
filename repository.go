package main

import (
	"fmt"
	"sync"
)

type StorageMemory struct {
	Users   []*User
	Tasks   []*Task
	MuUsers *sync.RWMutex
	MuTasks *sync.RWMutex
}

func NewStorage() *StorageMemory {
	return &StorageMemory{
		Tasks:   make([]*Task, 0, 10),
		Users:   make([]*User, 0, 10),
		MuUsers: &sync.RWMutex{},
		MuTasks: &sync.RWMutex{},
	}
}

func (t *StorageMemory) AddUser(userName string, chatID int64) error {

	t.MuUsers.Lock()
	t.Users = append(t.Users, &User{
		UserName: userName,
		ChatID:   chatID,
	})
	t.MuUsers.Unlock()
	return nil
}

func (t *StorageMemory) CheckUser(userName string) (*User, error) {
	t.MuUsers.RLock()
	defer t.MuUsers.RUnlock()

	for _, user := range t.Users {
		if user.UserName == userName {
			return user, nil
		}
	}
	return nil, fmt.Errorf("there is no user with this username")
}

func (t *StorageMemory) GetAllTasks() ([]*Task, error) {
	t.MuTasks.RLock()
	tasks := t.Tasks
	t.MuTasks.RUnlock()
	return tasks, nil
}

func (t *StorageMemory) AddTask(user, task string) (int, error) {
	t.MuTasks.RLock()
	id := len(t.Tasks) + 1
	t.MuTasks.RUnlock()

	new := &Task{
		ID:    id,
		Name:  task,
		Owner: user,
	}
	t.MuTasks.Lock()
	t.Tasks = append(t.Tasks, new)
	t.MuTasks.Unlock()

	return id, nil
}

func (t *StorageMemory) AddAsigner(user string, id int) error {
	t.MuTasks.Lock()
	defer t.MuTasks.Unlock()
	for i := 0; i < len(t.Tasks); i++ {
		if t.Tasks[i].ID == id {

			t.Tasks[i].Assigned = user

			return nil
		}
	}
	return fmt.Errorf("assigned is not updated, no task with this id")
}

func (t *StorageMemory) GetTaskWithID(id int) (*Task, error) {
	t.MuTasks.RLock()
	defer t.MuTasks.RUnlock()
	for i := 0; i < len(t.Tasks); i++ {
		if t.Tasks[i].ID == id {

			return t.Tasks[i], nil
		}
	}
	return nil, fmt.Errorf("no task with this id")
}

func (t *StorageMemory) SetDoneToTask(id int) error {
	t.MuTasks.RLock()
	defer t.MuTasks.RUnlock()
	for i := 0; i < len(t.Tasks); i++ {
		if t.Tasks[i].ID == id {
			t.Tasks[i].Done = true
			return nil
		}
	}
	return fmt.Errorf("no task with this id")
}
