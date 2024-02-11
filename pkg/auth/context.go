package auth

import (
	"htmx/model"
)

type UserInfo struct {
	User    model.User
	Session model.UserSession
}
