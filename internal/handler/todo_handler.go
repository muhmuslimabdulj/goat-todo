package handler

import (
	"net/http"
	"strconv"
	"strings"

	"todo-goat/internal/app"
	"todo-goat/internal/handler/templates"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	service *app.TodoService
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(service *app.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

// Index renders the main page
func (h *TodoHandler) Index(w http.ResponseWriter, r *http.Request) {
	todos, _ := h.service.GetAllTodos()
	templates.Index(todos).Render(r.Context(), w)
}

// Create creates a new todo and returns updated container
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	h.service.CreateTodo(title)
	todos, _ := h.service.GetAllTodos()
	templates.TodoContainer(todos).Render(r.Context(), w)
}

// Toggle toggles a todo and returns updated container
func (h *TodoHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractID(r.URL.Path)
	if id == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	h.service.ToggleTodo(id)
	todos, _ := h.service.GetAllTodos()
	templates.TodoContainer(todos).Render(r.Context(), w)
}

// Delete deletes a todo and returns updated container
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractID(r.URL.Path)
	if id == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	h.service.DeleteTodo(id)
	todos, _ := h.service.GetAllTodos()
	templates.TodoContainer(todos).Render(r.Context(), w)
}

// RegisterRoutes registers all todo routes
func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/todos", h.Create)
	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/toggle") {
			h.Toggle(w, r)
		} else if strings.HasSuffix(path, "/delete") {
			h.Delete(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}

// extractID extracts the ID from a path like /todos/123/toggle
func extractID(path string) int {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}
	return id
}
