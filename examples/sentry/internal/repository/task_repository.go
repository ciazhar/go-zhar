package repository

import (
	"database/sql"
	"errors"
	"github.com/ciazhar/go-start-small/examples/sentry/internal/model"

	"github.com/getsentry/sentry-go"
)

var db *sql.DB

// SetDB assigns the database instance to the repository
func SetDB(database *sql.DB) {
	db = database
}

// CreateTask inserts a new task into the database
func CreateTask(task model.Task) error {
	query := `INSERT INTO tasks (title, description) VALUES ($1, $2)`
	_, err := db.Exec(query, task.Title, task.Description)
	if err != nil {
		sentry.CaptureException(err)
		if isUniqueViolation(err) {
			return errors.New("task with this title already exists")
		}
		return err
	}
	return nil
}

// GetAllTasks fetches all tasks from the database
func GetAllTasks() ([]model.Task, error) {
	rows, err := db.Query("SELECT id, title, description FROM tasks")
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
			sentry.CaptureException(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Helper function to detect unique constraint violations
func isUniqueViolation(err error) bool {
	if err.Error() == `pq: duplicate key value violates unique constraint "tasks_title_key"` {
		return true
	}
	return false
}
