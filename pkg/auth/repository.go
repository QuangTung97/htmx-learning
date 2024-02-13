package auth

import (
	"context"

	"github.com/QuangTung97/svloc"
	"github.com/pkg/errors"

	"htmx/model"
	"htmx/pkg/dbtx"
	"htmx/pkg/util/dblib"
)

//go:generate moq -out repository_mocks_test.go . Repository RandService

var ErrDuplicatedUser = errors.New("duplicated user")

type Repository interface {
	GetUser(ctx context.Context, userID model.UserID) (model.NullUser, error)
	FindUserSession(ctx context.Context, userID model.UserID, sessionID model.SessionID) (model.NullUserSession, error)

	InsertUser(ctx context.Context, user model.User) (model.UserID, error)
	InsertUserSession(ctx context.Context, sess model.UserSession) error
}

type repoImpl struct {
}

var RepoLoc = svloc.Register[Repository](func(unv *svloc.Universe) Repository {
	return &repoImpl{}
})

func (r *repoImpl) GetUser(ctx context.Context, userID model.UserID) (model.NullUser, error) {
	query := `
SELECT id, provider, oauth_user_id, email, image_url, status
FROM users WHERE id = ?
`
	return dblib.Get[model.User](ctx, query, userID)
}

func (r *repoImpl) FindUserSession(
	ctx context.Context, userID model.UserID, sessionID model.SessionID,
) (model.NullUserSession, error) {
	query := `
SELECT user_id, session_id, status FROM user_sessions
WHERE user_id = ? AND session_id = ?
`
	return dblib.Get[model.UserSession](ctx, query, userID, sessionID)
}

func (r *repoImpl) InsertUser(ctx context.Context, user model.User) (model.UserID, error) {
	query := `
INSERT INTO users (provider, oauth_user_id, email, image_url, status)
VALUES (:provider, :oauth_user_id, :email, :image_url, :status)
`
	userID, err := dblib.Insert[model.UserID](ctx, query, user)
	if dblib.IsDuplicatedErr(err) {
		return 0, ErrDuplicatedUser
	}
	return userID, err
}

func (r *repoImpl) InsertUserSession(ctx context.Context, sess model.UserSession) error {
	query := `
INSERT INTO user_sessions (user_id, session_id, status)
VALUES (:user_id, :session_id, :status)
`
	tx := dbtx.GetTx(ctx)
	_, err := tx.NamedExecContext(ctx, query, sess)
	return err
}
