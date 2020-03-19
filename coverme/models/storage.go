// +build !change

package models

import (
	"fmt"
	"sync"
)

type Storage interface {
	AddTodo(string, string) (*Todo, error)
	GetTodo(ID) (*Todo, error)
	GetAll() ([]*Todo, error)
}

type InMemoryStorage struct {
	mu    sync.RWMutex
	todos map[ID]*Todo

	nextID ID
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		todos: make(map[ID]*Todo),
	}
}

func (s *InMemoryStorage) AddTodo(title, content string) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	todo := &Todo{
		ID:       id,
		Title:    title,
		Content:  content,
		Finished: false,
	}

	s.todos[todo.ID] = todo

	return todo, nil
}

func (s *InMemoryStorage) GetTodo(id ID) (*Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo %d not found", id)
	}

	return todo, nil
}

func (s *InMemoryStorage) GetAll() ([]*Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []*Todo
	for _, todo := range s.todos {
		out = append(out, todo)
	}

	return out, nil
}
