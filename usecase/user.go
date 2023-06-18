package usecase

import (
	"context"
	"errors"
	"tinder-like-app/container"
	"tinder-like-app/helper"
	"tinder-like-app/model"
	"tinder-like-app/repository"
)

type User struct {
	userRepo repository.User
}

func NewUser(app *container.Container) User {
	return User{
		userRepo: app.UserRepo(),
	}
}

func (u *User) SignUp(ctx context.Context, req model.User) (*model.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// check existing email
	user, err := u.userRepo.GetByEmail(ctx, *req.Email)
	if !helper.IsNotFoundError(err) || user != nil {
		return nil, errors.New("email already exist")
	}
	if err != nil && !helper.IsNotFoundError(err) {
		return nil, err
	}

	// add new user
	newUser, err := u.userRepo.Add(ctx, req)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *User) Login(ctx context.Context, req model.User) (*model.User, bool) {
	// get user by email
	user, err := u.userRepo.GetByEmail(ctx, *req.Email)
	if err != nil || user == nil {
		return nil, false
	}

	// validate password
	if *user.Password != helper.Hash(*user.PasswordSalt, *req.Password) {
		return nil, false
	}

	return user, true
}
