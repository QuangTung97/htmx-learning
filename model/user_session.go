package model

import (
	"time"

	"htmx/pkg/util"
)

type SessionID string

type UserSession struct {
	UserID    UserID    `db:"user_id"`
	SessionID SessionID `db:"session_id"`

	Status UserSessionStatus `db:"status"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type NullUserSession = util.Null[UserSession]

type UserSessionStatus int

const (
	UserSessionStatusActive UserSessionStatus = iota + 1
	UserSessionStatusDeleted
)
