package bun_entities

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`

	ID           int       `bun:"id,pk,autoincrement"`
	UserID       int       `bun:"user_id,notnull"`
	RefreshToken string    `bun:"refresh_token,notnull"`
	ExpiresAt    time.Time `bun:"expires_at,notnull"`
}
