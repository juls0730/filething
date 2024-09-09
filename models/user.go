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
	bun.BaseModel `bun:"table:users"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	Username      string    `bun:"username,notnull,unique" json:"username"`
	Email         string    `bun:"email,notnull,unique" json:"email"`
	PasswordHash  string    `bun:"passwordHash,notnull" json:"-"`
	PlanID        int64     `bun:"plan_id,notnull" json:"-"`
	Plan          Plan      `bun:"rel:belongs-to,join:plan_id=id" json:"plan"`
	Usage         int64     `bun:"-" json:"usage"`
}

type Session struct {
	bun.BaseModel `bun:"table:sessions"`
	ID            uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID        uuid.UUID `bun:"user_id,notnull,type:uuid"`
	User          User      `bun:"rel:belongs-to,join:user_id=id"`
}

type Plan struct {
	bun.BaseModel `bun:"table:plans"`
	ID            int64 `bun:"id,pk,autoincrement" json:"id"`
	MaxStorage    int64 `bun:"max_storage,notnull" json:"max_storage"`
}
