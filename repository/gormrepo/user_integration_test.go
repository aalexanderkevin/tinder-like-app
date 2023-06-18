//go:build integration
// +build integration

package gormrepo_test

import (
	"context"
	"testing"
	"tinder-like-app/helper"
	"tinder-like-app/helper/test"
	"tinder-like-app/model"
	"tinder-like-app/repository/gormrepo"
	"tinder-like-app/storage"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Add(t *testing.T) {
	t.Run("ShouldInsertUser", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUser(t, nil)

		//-- code under test
		userRepo := gormrepo.NewUserRepository(db)
		addedUser, err := userRepo.Add(context.TODO(), *fakeUser)

		//-- assert
		require.NoError(t, err)
		require.NotNil(t, addedUser)
		existingUser, err := userRepo.Get(context.TODO(), *addedUser.ID)
		require.NoError(t, err)
		require.NotNil(t, existingUser)
		require.Equal(t, addedUser.ID, existingUser.ID)
		require.Equal(t, addedUser.Name, existingUser.Name)
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("ShouldNotFoundError_WhenIdNotExist", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUser(t, nil)

		//-- code under test
		userRepo := gormrepo.NewUserRepository(db)
		resUpdate, err := userRepo.Update(context.Background(), "test-id", *fakeUser)

		//-- assert
		require.Error(t, err)
		require.Nil(t, resUpdate)
	})

	t.Run("ShouldUpdateUserPassword", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUserCreate(t, db, nil)
		fakeUser2 := model.User{
			Password:     helper.Pointer(fake.CharactersN(7)),
			PasswordSalt: helper.Pointer(fake.CharactersN(7)),
		}

		//-- code under test
		userRepo := gormrepo.NewUserRepository(db)
		resUpdate, err := userRepo.Update(context.Background(), *fakeUser.ID, fakeUser2)
		data, err := userRepo.Get(context.Background(), *resUpdate.ID)

		//-- assert
		require.NoError(t, err)
		require.Equal(t, *resUpdate.Password, *fakeUser2.Password)
		require.Equal(t, *resUpdate.PasswordSalt, *data.PasswordSalt)
		require.NotEqual(t, *fakeUser.Password, *fakeUser2.Password)
	})
}

func TestUserRepository_Get(t *testing.T) {
	t.Run("ShouldReturnError_WhenIDIsNotFound", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		//-- code under test
		userRepo := gormrepo.NewUserRepository(db)
		res, err := userRepo.Get(context.TODO(), "test")

		//-- assert
		require.Error(t, err)
		require.True(t, helper.IsNotFoundError(err))
		require.Nil(t, res)
	})

	t.Run("ShouldReturnUser_WhenIdIsFound", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUserCreate(t, db, nil)

		//-- code under test
		userRepo := gormrepo.NewUserRepository(db)
		res, err := userRepo.Get(context.TODO(), *fakeUser.ID)

		//-- assert
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, fakeUser.ID, res.ID)
		require.Equal(t, fakeUser.Name, res.Name)
	})
}

func TestUserRepository_GetEmail(t *testing.T) {
	t.Run("ShouldReturnError_WhenEmailIsNotFound", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		//-- code under test
		userRepo := gormrepo.NewUserRepository(db)
		res, err := userRepo.GetByEmail(context.TODO(), "test@test.com")

		//-- assert
		require.Error(t, err)
		require.True(t, helper.IsNotFoundError(err))
		require.Nil(t, res)
	})

	t.Run("ShouldReturnUser_WhenIdIsFound", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUserCreate(t, db, nil)

		//-- code under test
		userRepo := gormrepo.NewUserRepository(db)
		res, err := userRepo.GetByEmail(context.TODO(), *fakeUser.Email)

		//-- assert
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, fakeUser.ID, res.ID)
		require.Equal(t, fakeUser.Name, res.Name)
		require.Equal(t, fakeUser.Password, res.Password)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	t.Run("ShouldDeleteUser", func(t *testing.T) {
		//-- init
		db := storage.PostgresDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUserCreate(t, db, nil)
		test.FakeUserCreate(t, db, nil)

		//-- code under test
		userRepo := gormrepo.NewUserRepository(db)
		err := userRepo.Delete(context.Background(), *fakeUser.ID)
		require.NoError(t, err)
		data, err := userRepo.Get(context.Background(), *fakeUser.ID)
		require.Error(t, err)
		require.Nil(t, data)
	})
}
