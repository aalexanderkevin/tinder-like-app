package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	Id          *string `json:"id,omitempty" example:"user-id"`
	Name        *string `json:"name,omitempty" example:"name"`
	Email       *string `json:"email,omitempty" example:"email"`
	Password    *string `json:"password,omitempty" example:"password"`
	Description *string `json:"description"`
	Gender      *string `json:"gender"`
	IsPremium   bool    `json:"is_premium,omitempty" example:"10"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Gender, validation.Required, validation.In("f", "m")),
		validation.Field(&u.Description, validation.Length(1, 255)),
	)
}

func (u User) ValidateLogin() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)
}
