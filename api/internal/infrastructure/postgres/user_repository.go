package postgres

import (
	"context"
	"database/sql"
	"errors"
	"hilo-api/internal/domain/do"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *do.User) error {
	query := `
		INSERT INTO users (id, email, password, username, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID(),
		user.Email(),
		user.PasswordHash(),
		user.Username(),
		user.CreatedAt(),
	)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*do.User, error) {
	query := `
		SELECT id, email, password, username, created_at
		FROM users
		WHERE id = $1
	`

	var (
		uid       uuid.UUID
		email     string
		password  string
		username  string
		createdAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&uid, &email, &password, &username, &createdAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return do.ReconstructUser(uid, email, password, username, createdAt.Time), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*do.User, error) {
	query := `
		SELECT id, email, password, username, created_at
		FROM users
		WHERE email = $1
	`

	var (
		id        uuid.UUID
		userEmail string
		password  string
		username  string
		createdAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&id, &userEmail, &password, &username, &createdAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return do.ReconstructUser(id, userEmail, password, username, createdAt.Time), nil
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]*do.User, error) {
	query := `
		SELECT id, email, password, username, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*do.User
	for rows.Next() {
		var (
			id        uuid.UUID
			email     string
			password  string
			username  string
			createdAt sql.NullTime
		)

		if err := rows.Scan(&id, &email, &password, &username, &createdAt); err != nil {
			return nil, err
		}

		users = append(users, do.ReconstructUser(id, email, password, username, createdAt.Time))
	}

	return users, rows.Err()
}

func (r *UserRepository) Search(ctx context.Context, queryString string, limit int) ([]*do.User, error) {
	query := `
		SELECT id, email, password, username, created_at
		FROM users
		WHERE username ILIKE $1
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, "%"+queryString+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*do.User
	for rows.Next() {
		var (
			id        uuid.UUID
			email     string
			password  string
			username  string
			createdAt sql.NullTime
		)

		if err := rows.Scan(&id, &email, &password, &username, &createdAt); err != nil {
			return nil, err
		}

		users = append(users, do.ReconstructUser(id, email, password, username, createdAt.Time))
	}

	return users, rows.Err()
}
