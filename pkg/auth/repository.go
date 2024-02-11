package auth

import (
	"context"

	"htmx/model"
)

type Repository interface {
	GetUser(ctx context.Context, userID model.UserID) (model.NullUser, error)
	FindUserSession(ctx context.Context, sessionID model.SessionID) (model.NullUserSession, error)

	InsertUser(ctx context.Context, user model.User) (model.UserID, error)
	InsertUserSession(ctx context.Context, userID model.UserID, sessionID model.SessionID) error
}

type repoImpl struct {
}

func NewRepository() Repository {
	return &repoImpl{}
}

func (r *repoImpl) GetUser(ctx context.Context, userID model.UserID) (model.NullUser, error) {
	return model.NullUser{}, nil
}

func (r *repoImpl) FindUserSession(
	ctx context.Context, sessionID model.SessionID,
) (model.NullUserSession, error) {
	return model.NullUserSession{}, nil
}

func (r *repoImpl) InsertUser(ctx context.Context, user model.User) (model.UserID, error) {
	return 0, nil
}

func (r *repoImpl) InsertUserSession(ctx context.Context, userID model.UserID, sessionID model.SessionID) error {
	return nil
}
