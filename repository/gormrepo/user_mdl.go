package gormrepo

import (
	"time"
	"tinder-like-app/model"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type User struct {
	ID           *string    `json:"id"`
	Name         *string    `json:"name"`
	Email        *string    `json:"email"`
	PasswordSalt *string    `json:"password_salt"`
	Password     *string    `json:"no_share"`
	CreatedAt    *time.Time `json:"created_at"`
}

func (u User) GetID() *string {
	return u.ID
}

func (u User) TableName() string {
	return "users"
}

func (u User) FromModel(data model.User) *User {
	return &User{
		ID:           data.ID,
		Name:         data.Name,
		Email:        data.Email,
		PasswordSalt: data.PasswordSalt,
		Password:     data.Password,
		CreatedAt:    data.CreatedAt,
	}
}

func (u User) ToModel() *model.User {
	return &model.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordSalt: u.PasswordSalt,
		Password:     u.Password,
		CreatedAt:    u.CreatedAt,
	}
}

func (u User) ToModels(users []User) (ret []model.User) {
	for _, v := range users {
		m := v.ToModel()
		ret = append(ret, *m)
	}
	return ret
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.ID == nil {
		db.Statement.SetColumn("ID", ksuid.New().String())
	}
	return nil
}
