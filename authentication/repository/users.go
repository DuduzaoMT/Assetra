package repository

import (
	"assetra/authentication/models"
	"assetra/db"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const dbTimeout = 5 * time.Second

type UserRepository interface {
	// Define methods for user repository
	Save(user *models.User) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUserName(username string) (*models.User, error)
	FindAll() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	
	// Login attempt tracking
	IncrementFailedLoginAttempts(email string) error
	ResetFailedLoginAttempts(email string) error
	LockAccount(email string, until time.Time) error
	IsAccountLocked(email string) (bool, error)
}

type userRepository struct {
	// Using pgxpool for concurrent access
	database *pgxpool.Pool
}

func NewUserRepository(db db.Connection) UserRepository {
	return &userRepository{database: db.DB()}
}

func (r *userRepository) Save(user *models.User) (*models.User, error) {
	query := `
        INSERT INTO users (username, email, password, created_at, updated_at, roles) 
        VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
    `

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	err := r.database.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
		user.Roles,
	).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, roles, 
		       failed_login_attempts, locked_until
		FROM users WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user models.User
	err := r.database.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Roles,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, roles,
		       failed_login_attempts, locked_until
		FROM users WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user models.User
	err := r.database.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Roles,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUserName(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, roles,
		       failed_login_attempts, locked_until
		FROM users WHERE username = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user models.User
	err := r.database.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Roles,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, roles,
		       failed_login_attempts, locked_until
		FROM users
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := r.database.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.Password,
			&u.CreatedAt, &u.UpdatedAt, &u.Roles,
			&u.FailedLoginAttempts, &u.LockedUntil,
		); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return users, nil
}

func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users SET username = $1, email = $2, password = $3, roles = $4, updated_at = $5 WHERE id = $6
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := r.database.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Roles,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

func (r *userRepository) Delete(id string) error {
	query := `
		DELETE FROM users WHERE id = $1
	`
	
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	
	_, err := r.database.Exec(ctx, query, id)
	return err
}

func (r *userRepository) IncrementFailedLoginAttempts(email string) error {
	query := `
		UPDATE users 
		SET failed_login_attempts = failed_login_attempts + 1 
		WHERE email = $1
	`
	
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	
	_, err := r.database.Exec(ctx, query, email)
	return err
}

func (r *userRepository) ResetFailedLoginAttempts(email string) error {
	query := `
		UPDATE users 
		SET failed_login_attempts = 0, locked_until = NULL
		WHERE email = $1
	`
	
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	
	_, err := r.database.Exec(ctx, query, email)
	return err
}

func (r *userRepository) LockAccount(email string, until time.Time) error {
	query := `
		UPDATE users 
		SET locked_until = $1 
		WHERE email = $2
	`
	
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	
	_, err := r.database.Exec(ctx, query, until, email)
	return err
}

func (r *userRepository) IsAccountLocked(email string) (bool, error) {
	query := `
		SELECT COALESCE(locked_until, '1970-01-01'::timestamp) > NOW() as is_locked
		FROM users 
		WHERE email = $1
	`
	
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	
	var isLocked bool
	err := r.database.QueryRow(ctx, query, email).Scan(&isLocked)
	if err != nil {
		return false, err
	}
	return isLocked, nil
}
