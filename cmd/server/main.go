package main

import (
	"log"
	"net/http"

	"todo-goat/internal/app"
	"todo-goat/internal/handler"
	"todo-goat/internal/infra"
)

func main() {
	// Infrastructure layer - database
	repo, err := infra.NewSQLiteTodoRepo("todo.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer repo.Close()

	// Seed database with sample data if empty
	if err := repo.SeedIfEmpty(); err != nil {
		log.Println("Warning: Failed to seed database:", err)
	}

	// Application layer - services
	todoService := app.NewTodoService(repo)

	// Infrastructure layer - event hub for SSE
	eventHub := infra.NewEventHub()

	// Interface layer - handlers
	todoHandler := handler.NewTodoHandler(todoService, eventHub)

	// Wire routes
	mux := http.NewServeMux()
	todoHandler.RegisterRoutes(mux)

	// Start server
	log.Println("🐐 Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
