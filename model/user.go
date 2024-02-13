package model

import (
	"time"

	"htmx/pkg/util"
)

type UserID int64

type User struct {
	ID UserID `db:"id"`

	Provider    string `db:"provider"`
	OAuthUserID string `db:"oauth_user_id"`
	Email       string `db:"email"`
	ImageURL    string `db:"image_url"`

	Status UserStatus `db:"status"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type NullUser = util.Null[User]

type UserStatus int

const (
	UserStatusActive UserStatus = iota + 1
)
