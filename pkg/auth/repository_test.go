//go:build integration

package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"htmx/model"
	"htmx/pkg/integration"
)

func TestRepository(t *testing.T) {
	t.Run("user", func(t *testing.T) {
		tc := integration.NewTestCase()
		tc.TruncateTables(model.User{})
		r := RepoLoc.Get(tc.Unv)

		ctx := tc.Autocommit()

		nullUser, err := r.GetUser(ctx, 1)
		assert.Equal(t, nil, err)
		assert.Equal(t, false, nullUser.Valid)

		user := model.User{
			Provider:    "google",
			OAuthUserID: "1234",
			Email:       "quangtung@gmail.com",
			ImageURL:    "image-path",
			Status:      model.UserStatusActive,
		}

		// Do Insert
		userID, err := r.InsertUser(ctx, user)
		assert.Equal(t, nil, err)
		assert.Equal(t, model.UserID(1), userID)

		user.ID = userID

		// Get User
		nullUser, err = r.GetUser(ctx, userID)
		assert.Equal(t, nil, err)
		assert.Equal(t, model.NullUser{
			Valid: true,
			Data:  user,
		}, nullUser)

		// Get Not Found
		nullUser, err = r.GetUser(ctx, 2)
		assert.Equal(t, nil, err)
		assert.Equal(t, false, nullUser.Valid)

		// Do Insert Duplicated
		user.ID = 0
		userID, err = r.InsertUser(ctx, user)
		assert.Equal(t, ErrDuplicatedUser, err)
		assert.Equal(t, model.UserID(0), userID)
	})

	t.Run("user", func(t *testing.T) {
		tc := integration.NewTestCase()
		tc.TruncateTables(model.UserSession{})
		r := RepoLoc.Get(tc.Unv)

		ctx := tc.Autocommit()

		sess := model.UserSession{
			UserID:    21,
			SessionID: "sess01",
			Status:    model.UserSessionStatusActive,
		}
		err := r.InsertUserSession(ctx, sess)
		assert.Equal(t, nil, err)

		findSess, err := r.FindUserSession(ctx, sess.UserID, sess.SessionID)
		assert.Equal(t, nil, err)
		assert.Equal(t, model.NullUserSession{
			Valid: true,
			Data:  sess,
		}, findSess)

		// Not found
		findSess, err = r.FindUserSession(ctx, sess.UserID+1, sess.SessionID)
		assert.Equal(t, nil, err)
		assert.Equal(t, model.NullUserSession{}, findSess)

		// Not found
		findSess, err = r.FindUserSession(ctx, sess.UserID, sess.SessionID+"a")
		assert.Equal(t, nil, err)
		assert.Equal(t, model.NullUserSession{}, findSess)
	})
}
