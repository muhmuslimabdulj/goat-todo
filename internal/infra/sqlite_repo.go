package infra

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"todo-goat/internal/domain"
)

// SQLiteTodoRepo implements TodoRepository with SQLite
type SQLiteTodoRepo struct {
	db *sql.DB
}

// NewSQLiteTodoRepo creates a new SQLiteTodoRepo
func NewSQLiteTodoRepo(dbPath string) (*SQLiteTodoRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteTodoRepo{db: db}, nil
}

// Create creates a new todo
func (r *SQLiteTodoRepo) Create(title string) (*domain.Todo, error) {
	result, err := r.db.Exec("INSERT INTO todos (title) VALUES (?)", title)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &domain.Todo{ID: int(id), Title: title, Done: false}, nil
}

// GetAll returns all todos ordered by ID descending
func (r *SQLiteTodoRepo) GetAll() ([]domain.Todo, error) {
	rows, err := r.db.Query("SELECT id, title, done FROM todos ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []domain.Todo
	for rows.Next() {
		var t domain.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

// GetByID returns a todo by ID
func (r *SQLiteTodoRepo) GetByID(id int) (*domain.Todo, error) {
	var t domain.Todo
	err := r.db.QueryRow("SELECT id, title, done FROM todos WHERE id = ?", id).Scan(&t.ID, &t.Title, &t.Done)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Toggle toggles the done status and returns the updated todo
func (r *SQLiteTodoRepo) Toggle(id int) (*domain.Todo, error) {
	_, err := r.db.Exec("UPDATE todos SET done = NOT done WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// Delete deletes a todo
func (r *SQLiteTodoRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

// SeedIfEmpty populates the database with sample todos if empty
func (r *SQLiteTodoRepo) SeedIfEmpty() error {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		seedData := []struct {
			title string
			done  bool
		}{
			{"Selamat datang di GOAT Todo!", true},
			{"Tambahkan todo baru dengan form di atas", false},
			{"Klik checkbox untuk menandai selesai", false},
			{"Klik tombol hapus untuk menghapus todo", false},
			{"Deploy berhasil!", true},
		}

		for _, item := range seedData {
			_, err := r.db.Exec("INSERT INTO todos (title, done) VALUES (?, ?)", item.title, item.done)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Close closes the database connection
func (r *SQLiteTodoRepo) Close() error {
	return r.db.Close()
}
