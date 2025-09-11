package user

const (
	queryCreateUser = `
		INSERT INTO users (username, email, password, full_name)
		VALUES ($1, $2, $3, $4)
	`

	queryGetUserByID = `
		SELECT id, username, email, full_name, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL;
	`

	queryGetUsersWithPagination = `
		SELECT id, username, email, full_name, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	queryCountUsers = `
		SELECT COUNT(*)
		FROM users
		WHERE deleted_at IS NULL;
	`

	queryIsUserExistsByEmail = `
		SELECT EXISTS (
			SELECT 1 
			FROM users
			WHERE email = $1 AND deleted_at IS NULL 
			LIMIT 1; 
		)
	`

	querySoftDeleteUser = `
		UPDATE users
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL;
	`

	queryUpdateUser = `
		UPDATE users
		SET username   = $1,
			email      = $2,
			full_name  = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND deleted_at IS NULL;
	`

	queryUpsertUser = `
		INSERT INTO users (id, username, email, password, full_name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id)
		DO UPDATE SET
			username   = EXCLUDED.username,
			email      = EXCLUDED.email,
			password   = EXCLUDED.password,
			full_name  = EXCLUDED.full_name,
			updated_at = CURRENT_TIMESTAMP
		WHERE users.deleted_at IS NULL;
	`
)
