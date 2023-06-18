package repository

import (
	"context"
	"tinder-like-app/model"
)

type User interface {
	Add(ctx context.Context, user model.User) (*model.User, error)
	Get(ctx context.Context, userID string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, userID string, notification model.User) (*model.User, error)
	Delete(ctx context.Context, userID string) error
}
