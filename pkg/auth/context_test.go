package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"htmx/model"
)

func TestUserInfo(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctx := context.Background()

		info, ok := GetUserInfoNull(ctx)
		assert.Equal(t, false, ok)
		assert.Equal(t, UserInfo{}, info)

		newInfo := UserInfo{
			User: model.User{
				ID:       12,
				Provider: "google",
				Email:    "test@gmail.com",
			},
			Session: model.UserSession{
				UserID:    12,
				SessionID: "session-num",
			},
		}
		ctx = SetUserInfo(ctx, newInfo)

		info, ok = GetUserInfoNull(ctx)
		assert.Equal(t, true, ok)
		assert.Equal(t, newInfo, info)

		info = GetUserInfo(ctx)
		assert.Equal(t, newInfo, info)
	})

	t.Run("panic when not found", func(t *testing.T) {
		ctx := context.Background()
		assert.PanicsWithValue(t, "Missing UserInfo in context object", func() {
			GetUserInfo(ctx)
		})
	})
}
