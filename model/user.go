package model

import (
	"time"
	"tinder-like-app/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	ID           *string    `json:"id"`
	Name         *string    `json:"name"`
	Email        *string    `json:"email"`
	PasswordSalt *string    `json:"password_salt"`
	Password     *string    `json:"no_share"`
	CreatedAt    *time.Time `json:"created_at"`
}

func (b User) Validate() (err error) {
	if err = validation.ValidateStruct(&b,
		validation.Field(&b.Name, validation.Required),
		validation.Field(&b.Email, validation.Required),
		validation.Field(&b.Password, validation.Required),
	); err != nil {
		return helper.NewParameterError(helper.Pointer(err.Error()))
	}
	return nil
}
