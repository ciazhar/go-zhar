package repository

import (
	"database/sql"
	"github.com/ciazhar/go-zhar/examples/mysql/transactional/model"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Create(user *model.User) (int64, error) {
	result, err := r.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MySQLUserRepository) WithTransaction(fn func(tx *sql.Tx) error) error {
	// Start the transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Execute the passed function
	err = fn(tx)
	if err != nil {
		// If the function returns an error, rollback the transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Transaction rollback error: %v", rbErr)
		}
		return err
	}

	// Commit the transaction if everything went fine
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
