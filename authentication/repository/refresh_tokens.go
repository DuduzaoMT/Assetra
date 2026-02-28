package repository

import (
	"assetra/authentication/models"
	"assetra/db"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type RefreshTokenRepository interface {
	Save(token *models.RefreshToken) error
	FindByTokenHash(tokenHash string) (*models.RefreshToken, error)
	RevokeToken(tokenHash string) error
	RevokeAllUserTokens(userID string) error
	DeleteExpired() error
}

type refreshTokenRepository struct {
	database *pgxpool.Pool
}

func NewRefreshTokenRepository(db db.Connection) RefreshTokenRepository {
	return &refreshTokenRepository{database: db.DB()}
}

func (r *refreshTokenRepository) Save(token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at, revoked)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	err := r.database.QueryRow(ctx, query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	).Scan(&token.ID)

	return err
}

func (r *refreshTokenRepository) FindByTokenHash(tokenHash string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked = FALSE AND expires_at > NOW()
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var token models.RefreshToken
	err := r.database.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.Revoked,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *refreshTokenRepository) RevokeToken(tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE token_hash = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := r.database.Exec(ctx, query, tokenHash)
	return err
}

func (r *refreshTokenRepository) RevokeAllUserTokens(userID string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE user_id = $1 AND revoked = FALSE
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := r.database.Exec(ctx, query, userID)
	return err
}

func (r *refreshTokenRepository) DeleteExpired() error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW() OR revoked = TRUE
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := r.database.Exec(ctx, query)
	return err
}
