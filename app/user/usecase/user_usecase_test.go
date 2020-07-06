package usecase_test

import (
	"context"
	"github.com/fajardm/ewallet-example/app/base"
	"github.com/fajardm/ewallet-example/app/errorcode"
	"github.com/fajardm/ewallet-example/app/user/mocks"
	"github.com/fajardm/ewallet-example/app/user/model"
	"github.com/fajardm/ewallet-example/app/user/usecase"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const contextTimeout = 2 * time.Second

func prepareMockUser() model.User {
	return model.User{
		Model: base.Model{
			ID:        uuid.NewV4(),
			CreatedBy: uuid.NewV4(),
			CreatedAt: time.Now(),
		},
		Username:       "john",
		Email:          "john@gmail.com",
		MobilePhone:    "08199999999",
		HashedPassword: []byte("secret"),
	}
}

func TestStore(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockUser := prepareMockUser()

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByUsernameOrEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil, errorcode.ErrNotFound).Once()
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("model.User")).Return(nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo, contextTimeout)
		err := uc.Store(context.TODO(), mockUser)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("existing-name", func(t *testing.T) {
		mockUserRepo.On("GetByUsernameOrEmail", mock.Anything, mock.Anything, mock.Anything).Return(&mockUser, nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo, contextTimeout)
		err := uc.Store(context.TODO(), mockUser)
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockUser := prepareMockUser()

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.Anything, mock.Anything).Return(&mockUser, nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo, contextTimeout)
		res, err := uc.GetByID(context.TODO(), mockUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-not-found", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.Anything, mock.Anything).Return(nil, errorcode.ErrNotFound).Once()
		uc := usecase.NewUserUsecase(mockUserRepo, contextTimeout)
		res, err := uc.GetByID(context.TODO(), uuid.NewV4())
		assert.Error(t, err)
		assert.Nil(t, res)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockUser := prepareMockUser()

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.Anything, mock.Anything).Return(&mockUser, nil).Once()
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("model.User")).Return(nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo, contextTimeout)
		err := uc.Update(context.TODO(), mockUser)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-not-found", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.Anything, mock.Anything).Return(nil, errorcode.ErrNotFound).Once()
		uc := usecase.NewUserUsecase(mockUserRepo, contextTimeout)
		err := uc.Update(context.TODO(), mockUser)
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockUserRepo := new(mocks.Repository)
	mockUser := prepareMockUser()

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.Anything, mock.Anything).Return(&mockUser, nil).Once()
		mockUserRepo.On("Delete", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
		uc := usecase.NewUserUsecase(mockUserRepo, contextTimeout)
		err := uc.Delete(context.TODO(), mockUser.ID)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-not-found", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.Anything, mock.Anything).Return(nil, errorcode.ErrNotFound).Once()
		uc := usecase.NewUserUsecase(mockUserRepo, contextTimeout)
		err := uc.Delete(context.TODO(), uuid.NewV4())
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}
