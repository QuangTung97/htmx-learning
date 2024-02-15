package auth

import (
	"context"

	"htmx/model"
)

type UserInfo struct {
	User    model.User
	Session model.UserSession
}

type userCtxKeyType struct {
}

var userCtxKey = &userCtxKeyType{}

func GetUserInfoNull(ctx context.Context) (UserInfo, bool) {
	info, ok := ctx.Value(userCtxKey).(UserInfo)
	return info, ok
}

func GetUserInfo(ctx context.Context) UserInfo {
	info, ok := GetUserInfoNull(ctx)
	if !ok {
		panic("Missing UserInfo in context object")
	}
	return info
}

func SetUserInfo(ctx context.Context, info UserInfo) context.Context {
	return context.WithValue(ctx, userCtxKey, info)
}
