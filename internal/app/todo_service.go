package app

import (
	"errors"

	"todo-goat/internal/domain"
)

// TodoService handles todo business logic
type TodoService struct {
	repo domain.TodoRepository
}

// NewTodoService creates a new TodoService
func NewTodoService(repo domain.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

// CreateTodo creates a new todo with validation
func (s *TodoService) CreateTodo(title string) (*domain.Todo, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	return s.repo.Create(title)
}

// GetAllTodos returns all todos
func (s *TodoService) GetAllTodos() ([]domain.Todo, error) {
	return s.repo.GetAll()
}

// ToggleTodo toggles the done status of a todo
func (s *TodoService) ToggleTodo(id int) (*domain.Todo, error) {
	return s.repo.Toggle(id)
}

// DeleteTodo deletes a todo
func (s *TodoService) DeleteTodo(id int) error {
	return s.repo.Delete(id)
}

// CountDone counts completed todos
func (s *TodoService) CountDone(todos []domain.Todo) int {
	count := 0
	for _, t := range todos {
		if t.Done {
			count++
		}
	}
	return count
}
