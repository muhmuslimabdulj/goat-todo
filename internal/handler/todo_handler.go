package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"todo-goat/internal/app"
	"todo-goat/internal/handler/templates"
	"todo-goat/internal/infra"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	service *app.TodoService
	hub     *infra.EventHub
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(service *app.TodoService, hub *infra.EventHub) *TodoHandler {
	return &TodoHandler{service: service, hub: hub}
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
	h.hub.Broadcast("update")
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
	h.hub.Broadcast("update")
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
	h.hub.Broadcast("update")
	todos, _ := h.service.GetAllTodos()
	templates.TodoContainer(todos).Render(r.Context(), w)
}

// RegisterRoutes registers all todo routes
func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/todos", h.Create)
	mux.HandleFunc("/events", h.SSE)
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

// SSE handles Server-Sent Events for real-time updates
func (h *TodoHandler) SSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create client channel
	client := make(chan string, 10)
	h.hub.Register(client)
	defer h.hub.Unregister(client)

	// Get flusher for streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Send initial ping
	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	// Stream events
	for {
		select {
		case msg := <-client:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
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
