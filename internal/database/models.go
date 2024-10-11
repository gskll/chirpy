// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshToken struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword string
}
