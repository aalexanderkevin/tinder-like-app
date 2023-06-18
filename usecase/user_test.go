package usecase_test

import (
	"context"
	"testing"
	"tinder-like-app/container"
	"tinder-like-app/helper"
	"tinder-like-app/helper/test"
	"tinder-like-app/model"
	"tinder-like-app/repository/mocks"
	"tinder-like-app/usecase"

	"github.com/icrowley/fake"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUser_SignUp(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnError_WhenPasswordIsMissing", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}

		fakeUser := &model.User{}

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.SignUp(context.Background(), *fakeUser)
		require.Error(t, err)
		require.True(t, helper.IsParameterError(err))
		require.Nil(t, res)

	})

	t.Run("ShouldReturnError_WhenThereIsAnExistingUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}
		userMock := &mocks.User{}
		appContainer.SetUserRepo(userMock)

		fakeUser := test.FakeUser(t, nil)

		userMock.On("GetByEmail", mock.Anything, *fakeUser.Email).Return(fakeUser, nil).Once()

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.SignUp(context.Background(), *fakeUser)
		require.Error(t, err)
		require.EqualError(t, err, "email already exist")
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldAddNewUser_WhenEmailNotFound", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}
		userMock := &mocks.User{}
		appContainer.SetUserRepo(userMock)

		fakeUser := test.FakeUser(t, nil)

		userMock.On("GetByEmail", mock.Anything, *fakeUser.Email).Return(nil, helper.NewNotFoundError()).Once()
		userMock.On("Add", mock.Anything, *fakeUser).Return(fakeUser, nil).Once()

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, err := uc.SignUp(context.Background(), *fakeUser)
		require.NoError(t, err)
		require.Equal(t, *fakeUser.Email, *res.Email)

		userMock.AssertExpectations(t)
	})
}
func TestUser_Login(t *testing.T) {
	t.Parallel()

	t.Run("ShouldReturnFalse_WhenEmailNotFound", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}
		userMock := &mocks.User{}
		appContainer.SetUserRepo(userMock)

		fakeUser := test.FakeUser(t, nil)

		userMock.On("GetByEmail", mock.Anything, *fakeUser.Email).Return(nil, helper.NewNotFoundError()).Once()

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, isSuccess := uc.Login(context.Background(), *fakeUser)
		require.False(t, isSuccess)
		require.Nil(t, res)

		userMock.AssertExpectations(t)
	})

	t.Run("ShouldReturnTrue_WhenTheLoginIsSuccess", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}
		userMock := &mocks.User{}
		appContainer.SetUserRepo(userMock)

		password := fake.SimplePassword()
		fakeUser := test.FakeUser(t, func(m model.User) model.User {
			salt := ksuid.New().String()
			m.PasswordSalt = helper.Pointer(salt)
			m.Password = helper.Pointer(helper.Hash(salt, password))
			return m
		})

		userMock.On("GetByEmail", mock.Anything, *fakeUser.Email).Return(fakeUser, nil).Once()

		// CODE UNDER TEST
		uc := usecase.NewUser(&appContainer)
		res, isSuccess := uc.Login(context.Background(), model.User{
			Email:    fakeUser.Email,
			Password: &password,
		})
		require.True(t, isSuccess)
		require.Equal(t, *fakeUser.Email, *res.Email)

		userMock.AssertExpectations(t)
	})
}
