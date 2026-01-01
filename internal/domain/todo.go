package domain

// Todo represents a todo item entity
type Todo struct {
	ID    int
	Title string
	Done  bool
}

// TodoRepository defines the interface for todo persistence
type TodoRepository interface {
	Create(title string) (*Todo, error)
	GetAll() ([]Todo, error)
	GetByID(id int) (*Todo, error)
	Toggle(id int) (*Todo, error)
	Delete(id int) error
}
