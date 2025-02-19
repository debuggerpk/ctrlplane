// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chat_links.sql

package entities

import (
	"context"

	"github.com/google/uuid"
)

const createChatLink = `-- name: CreateChatLink :one
INSERT INTO chat_links (hook, kind, link_to, data)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at, hook, kind, link_to, data
`

type CreateChatLinkParams struct {
	Hook   int32     `json:"hook"`
	Kind   string    `json:"kind"`
	LinkTo uuid.UUID `json:"link_to"`
	Data   []byte    `json:"data"`
}

func (q *Queries) CreateChatLink(ctx context.Context, arg CreateChatLinkParams) (ChatLink, error) {
	row := q.db.QueryRow(ctx, createChatLink,
		arg.Hook,
		arg.Kind,
		arg.LinkTo,
		arg.Data,
	)
	var i ChatLink
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Hook,
		&i.Kind,
		&i.LinkTo,
		&i.Data,
	)
	return i, err
}

const getChatLink = `-- name: GetChatLink :one
SELECT id, created_at, updated_at, hook, kind, link_to, data
FROM chat_links
WHERE link_to = $1
`

func (q *Queries) GetChatLink(ctx context.Context, linkTo uuid.UUID) (ChatLink, error) {
	row := q.db.QueryRow(ctx, getChatLink, linkTo)
	var i ChatLink
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Hook,
		&i.Kind,
		&i.LinkTo,
		&i.Data,
	)
	return i, err
}
