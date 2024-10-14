package repository

import (
	"database/sql"
	"github.com/ciazhar/go-zhar/examples/mysql/crud/model"
	_ "github.com/go-sql-driver/mysql"
)

type UserRepository interface {
	Create(user *model.User) (int64, error)
	GetByID(id int64) (*model.User, error)
	Update(user *model.User) error
	Delete(id int64) error
	FindAll(name, email *string, limit, offset int) ([]*model.User, error)
	Count(name, email *string) (int64, error)
	FindWithCursor(lastID *int64, name, email *string, limit int) ([]*model.User, error)
}

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (repo *MySQLUserRepository) Create(user *model.User) (int64, error) {
	result, err := repo.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (repo *MySQLUserRepository) GetByID(id int64) (*model.User, error) {
	user := &model.User{}
	err := repo.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *MySQLUserRepository) Update(user *model.User) error {
	_, err := repo.db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

func (repo *MySQLUserRepository) Delete(id int64) error {
	_, err := repo.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (repo *MySQLUserRepository) FindAll(name, email *string, limit, offset int) ([]*model.User, error) {
	var users []*model.User
	var args []interface{}
	query := "SELECT id, name, email FROM users WHERE 1=1"

	if name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*name+"%")
	}

	if email != nil {
		query += " AND email LIKE ?"
		args = append(args, "%"+*email+"%")
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *MySQLUserRepository) Count(name, email *string) (int64, error) {
	var args []interface{}
	query := "SELECT COUNT(*) FROM users WHERE 1=1"

	if name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*name+"%")
	}

	if email != nil {
		query += " AND email LIKE ?"
		args = append(args, "%"+*email+"%")
	}

	var count int64
	err := repo.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *MySQLUserRepository) FindWithCursor(lastID *int64, name, email *string, limit int) ([]*model.User, error) {
	var users []*model.User
	var args []interface{}
	query := "SELECT id, name, email FROM users WHERE 1=1"

	if lastID != nil {
		query += " AND id > ?"
		args = append(args, *lastID)
	}

	if name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*name+"%")
	}

	if email != nil {
		query += " AND email LIKE ?"
		args = append(args, "%"+*email+"%")
	}

	query += " ORDER BY id ASC LIMIT ?"
	args = append(args, limit)

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
