package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type LoginData struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type SignupData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Username      string    `bun:"username,notnull,unique"`
	Email         string    `bun:"email,notnull,unique"`
	PasswordHash  string    `bun:"passwordHash,notnull"`
}

type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:u"`
	ID            uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	UserID        uuid.UUID `bun:"user_id,notnull"`
	User          User      `bun:"rel:belongs-to,join:user_id=id"`
}
