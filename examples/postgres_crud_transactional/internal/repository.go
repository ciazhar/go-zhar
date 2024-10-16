package internal

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
    Create(ctx context.Context, user *User) (int, error)
    GetByID(ctx context.Context, id int) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
    BatchCreateWithTrx(ctx context.Context, users []User) error
}

type PgxUserRepository struct {
    pool *pgxpool.Pool
}

func NewPgxUserRepository(pool *pgxpool.Pool) *PgxUserRepository {
	return &PgxUserRepository{
		pool: pool,
	}
}

// Create inserts a new user into the database
func (r *PgxUserRepository) Create(ctx context.Context, user *User) (int, error) {
    var id int
    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
    err := r.pool.QueryRow(ctx, query, user.Name, user.Email).Scan(&id)
    if err != nil {
        return 0, fmt.Errorf("error creating user: %w", err)
    }
    return id, nil
}

// GetByID fetches a user by ID
func (r *PgxUserRepository) GetByID(ctx context.Context, id int) (*User, error) {
    user := &User{}
    query := `SELECT id, name, email FROM users WHERE id = $1`
    err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, fmt.Errorf("error fetching user: %w", err)
    }
    return user, nil
}

// GetAll fetches all users from the database
func (r *PgxUserRepository) GetAll(ctx context.Context) ([]User, error) {
	users := []User{}
	query := `SELECT id, name, email FROM users`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}
	return users, nil
}

// Update modifies an existing user in the database
func (r *PgxUserRepository) Update(ctx context.Context, user *User) error {
    query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
    _, err := r.pool.Exec(ctx, query, user.Name, user.Email, user.ID)
    if err != nil {
        return fmt.Errorf("error updating user: %w", err)
    }
    return nil
}

// Delete removes a user by ID from the database
func (r *PgxUserRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("error deleting user: %w", err)
    }
    return nil
}

// BatchCreateWithTrx inserts multiple users in a transaction
func (r *PgxUserRepository) BatchCreateWithTrx(ctx context.Context, users []User) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx) // Rollback in case anything goes wrong

    query := `INSERT INTO users (name, email) VALUES ($1, $2)`
    for _, user := range users {
        _, err := tx.Exec(ctx, query, user.Name, user.Email)
        if err != nil {
            return fmt.Errorf("error during batch insert: %w", err)
        }
    }

    err = tx.Commit(ctx) // Commit transaction if all went well
    if err != nil {
        return fmt.Errorf("error committing transaction: %w", err)
    }

    return nil
}