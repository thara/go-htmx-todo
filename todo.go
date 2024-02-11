package main

import (
	"sync"

	"github.com/gofrs/uuid"
)

type Task struct {
	ID   string
	Note string
	Done bool
}

type Store struct {
	mu   sync.RWMutex
	todo []Task
}

func (s *Store) GetTodoList() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Task(nil), s.todo...)
}

func (s *Store) GetItemLeft() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var count int
	for _, n := range s.todo {
		if !n.Done {
			count++
		}
	}
	return count
}

func (s *Store) AddTask(note string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := Task{
		ID:   uuid.Must(uuid.NewV4()).String(),
		Note: note,
		Done: false,
	}
	s.todo = append(s.todo, n)
	return n
}

func (s *Store) FindToDoByID(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.todo {
		if s.todo[i].ID == id {
			return s.todo[i], true
		}
	}
	return Task{}, false
}

func (s *Store) Update(id string, note string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todo {
		if s.todo[i].ID == id {
			s.todo[i].Note = note
			return true
		}
	}
	return false
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todo {
		if s.todo[i].ID == id {
			s.todo = append(s.todo[:i], s.todo[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) UpdateStatus(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todo {
		if s.todo[i].ID == id {
			s.todo[i].Done = !s.todo[i].Done
			return true
		}
	}
	return false
}

func (s *Store) ClearCompleted() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var newTodo []Task
	for _, n := range s.todo {
		if !n.Done {
			newTodo = append(newTodo, n)
		}
	}
	s.todo = newTodo
}
