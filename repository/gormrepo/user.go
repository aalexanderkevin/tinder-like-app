package gormrepo

import (
	"context"
	"errors"
	"tinder-like-app/helper"
	"tinder-like-app/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (u UserRepository) Get(ctx context.Context, userID string) (*model.User, error) {
	var err error

	userGorm := User{
		ID: &userID,
	}

	err = u.db.WithContext(ctx).First(&userGorm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError()
		}
		return nil, err
	}

	return userGorm.ToModel(), nil
}

func (u UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var err error
	userGorm := User{
		Email: &email,
	}

	err = u.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).First(&userGorm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError()
		}
		return nil, err
	}

	return userGorm.ToModel(), nil
}

func (u UserRepository) Add(ctx context.Context, user model.User) (ret *model.User, err error) {
	gormModel := User{}.FromModel(user)
	if err = u.db.WithContext(ctx).Create(&gormModel).Error; err != nil {
		return nil, err
	}

	return gormModel.ToModel(), nil
}

func (u UserRepository) Update(ctx context.Context, id string, user model.User) (resp *model.User, err error) {
	existingUser, err := u.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	newUser := User{}.FromModel(user)
	if err = u.db.WithContext(ctx).Model(&User{ID: existingUser.ID}).Updates(&User{
		Name:         newUser.Name,
		PasswordSalt: newUser.PasswordSalt,
		Password:     newUser.Password,
	}).Error; err != nil {
		return nil, err
	}

	return u.Get(ctx, *existingUser.ID)
}

func (u UserRepository) Delete(ctx context.Context, UserId string) error {
	return u.db.WithContext(ctx).Delete(User{}, "id = ?", UserId).Error
}
