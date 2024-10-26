// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: team_users.sql

package entities

import (
	"context"

	"github.com/google/uuid"
)

const createTeamUser = `-- name: CreateTeamUser :one
INSERT INTO team_users (team_id, user_id, role, is_active, is_admin)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, team_id, user_id, role, is_active, is_admin
`

type CreateTeamUserParams struct {
	TeamID   uuid.UUID `json:"team_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     TeamRole  `json:"role"`
	IsActive bool      `json:"is_active"`
	IsAdmin  bool      `json:"is_admin"`
}

func (q *Queries) CreateTeamUser(ctx context.Context, arg CreateTeamUserParams) (TeamUser, error) {
	row := q.db.QueryRow(ctx, createTeamUser,
		arg.TeamID,
		arg.UserID,
		arg.Role,
		arg.IsActive,
		arg.IsAdmin,
	)
	var i TeamUser
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TeamID,
		&i.UserID,
		&i.Role,
		&i.IsActive,
		&i.IsAdmin,
	)
	return i, err
}

const getTeamUser = `-- name: GetTeamUser :one
SELECT id, created_at, updated_at, team_id, user_id, role, is_active, is_admin
FROM team_users
WHERE user_id = $1
`

func (q *Queries) GetTeamUser(ctx context.Context, userID uuid.UUID) (TeamUser, error) {
	row := q.db.QueryRow(ctx, getTeamUser, userID)
	var i TeamUser
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TeamID,
		&i.UserID,
		&i.Role,
		&i.IsActive,
		&i.IsAdmin,
	)
	return i, err
}
