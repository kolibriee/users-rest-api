package bun_entities

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           int       `bun:"id,pk,autoincrement"`
	Role         string    `bun:"role,notnull"`
	Name         string    `bun:"name,notnull"`
	Username     string    `bun:"username,notnull"`
	Email        string    `bun:"email"`
	PasswordHash string    `bun:"password_hash,notnull"`
	City         string    `bun:"city,notnull"`
	RegisteredAt time.Time `bun:"registered_at,notnull,default:current_timestamp"`
}
