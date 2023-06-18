package test

import (
	"context"
	"testing"
	"tinder-like-app/helper"
	"tinder-like-app/model"
	"tinder-like-app/repository/gormrepo"

	"github.com/icrowley/fake"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func FakeUser(t *testing.T, cb func(m model.User) model.User) *model.User {
	t.Helper()
	salt := ksuid.New().String()
	pwd := helper.Hash(salt, fake.SimplePassword())

	user := model.User{
		ID:           helper.Pointer(fake.CharactersN(5)),
		Name:         helper.Pointer(fake.FullName()),
		Email:        helper.Pointer(fake.EmailAddress()),
		Password:     helper.Pointer(pwd),
		PasswordSalt: helper.Pointer(salt),
	}

	if cb != nil {
		user = cb(user)
	}

	return &user
}

func FakeUserCreate(t *testing.T, db *gorm.DB, cb func(m model.User) model.User) *model.User {
	t.Helper()

	fakeUser := FakeUser(t, cb)

	tenantUserRepo := gormrepo.NewUserRepository(db)
	_, err := tenantUserRepo.Add(context.Background(), *fakeUser)
	require.NoError(t, err)

	return fakeUser
}
