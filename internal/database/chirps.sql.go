// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: chirps.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, user_id, body, created_at, updated_at
`

type CreateChirpParams struct {
	Body   string
	UserID uuid.UUID
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.Body, arg.UserID)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteChirp = `-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id=$1
`

func (q *Queries) DeleteChirp(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteChirp, id)
	return err
}

const getChirp = `-- name: GetChirp :one
SELECT id, user_id, body, created_at, updated_at FROM chirps
WHERE id=$1
`

func (q *Queries) GetChirp(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getChirps = `-- name: GetChirps :many
SELECT id, user_id, body, created_at, updated_at FROM chirps
ORDER BY created_at ASC
`

func (q *Queries) GetChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, getChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Body,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
