// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: oauth_accounts.sql

package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createOAuthAccount = `-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (user_id, hook, hook_account_id, expires_at, type)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, hook, hook_account_id, expires_at, type
`

type CreateOAuthAccountParams struct {
	UserID        uuid.UUID   `json:"user_id"`
	Hook          string      `json:"hook"`
	HookAccountID string      `json:"hook_account_id"`
	ExpiresAt     time.Time   `json:"expires_at"`
	Type          pgtype.Text `json:"type"`
}

func (q *Queries) CreateOAuthAccount(ctx context.Context, arg CreateOAuthAccountParams) (OauthAccount, error) {
	row := q.db.QueryRow(ctx, createOAuthAccount,
		arg.UserID,
		arg.Hook,
		arg.HookAccountID,
		arg.ExpiresAt,
		arg.Type,
	)
	var i OauthAccount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Hook,
		&i.HookAccountID,
		&i.ExpiresAt,
		&i.Type,
	)
	return i, err
}

const getOAuthAccountByHookAccountID = `-- name: GetOAuthAccountByHookAccountID :one
SELECT id, created_at, updated_at, user_id, hook, hook_account_id, expires_at, type
FROM oauth_accounts
WHERE hook_account_id = $1 and hook = $2
`

type GetOAuthAccountByHookAccountIDParams struct {
	HookAccountID string `json:"hook_account_id"`
	Hook          string `json:"hook"`
}

func (q *Queries) GetOAuthAccountByHookAccountID(ctx context.Context, arg GetOAuthAccountByHookAccountIDParams) (OauthAccount, error) {
	row := q.db.QueryRow(ctx, getOAuthAccountByHookAccountID, arg.HookAccountID, arg.Hook)
	var i OauthAccount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Hook,
		&i.HookAccountID,
		&i.ExpiresAt,
		&i.Type,
	)
	return i, err
}

const getOAuthAccountByID = `-- name: GetOAuthAccountByID :one
SELECT id, created_at, updated_at, user_id, hook, hook_account_id, expires_at, type
FROM oauth_accounts
WHERE id = $1
`

func (q *Queries) GetOAuthAccountByID(ctx context.Context, id uuid.UUID) (OauthAccount, error) {
	row := q.db.QueryRow(ctx, getOAuthAccountByID, id)
	var i OauthAccount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Hook,
		&i.HookAccountID,
		&i.ExpiresAt,
		&i.Type,
	)
	return i, err
}

const getOAuthAccountsByUserID = `-- name: GetOAuthAccountsByUserID :many
SELECT id, created_at, updated_at, user_id, hook, hook_account_id, expires_at, type
FROM oauth_accounts
WHERE user_id = $1
`

func (q *Queries) GetOAuthAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]OauthAccount, error) {
	rows, err := q.db.Query(ctx, getOAuthAccountsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OauthAccount
	for rows.Next() {
		var i OauthAccount
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Hook,
			&i.HookAccountID,
			&i.ExpiresAt,
			&i.Type,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
