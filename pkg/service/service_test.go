package service

import (
	test_user "test-user"
	"test-user/pkg/repository"
	"testing"
	"time"
)

func TestUserService_CreateUser(t *testing.T) {
	service := NewUserService(repository.NewMemoryRepository())

	currentTS := time.Now().UnixNano()
	time.Sleep(time.Nanosecond * 100)

	user, err := service.CreateUser(test_user.User{
		Name: "My new user",
	})
	if err != nil {
		t.Error(err)
		return
	}

	if user.ID == "" {
		t.Error("ID is empty")
		return
	}

	if user.CreatedAt.UnixNano() < currentTS || user.UpdatedAt.UnixNano() < currentTS {
		t.Error("Time not set")
		return
	}

}

func TestUserService_CreateUser2(t *testing.T) {

	service := NewUserService(repository.NewMemoryRepository())

	_, err := service.CreateUser(test_user.User{})
	if _, ok := err.(*UserNameEmptyError); !ok {
		t.Error("Not check: name required")
		return
	}
}

func TestUserService_UpdateUser(t *testing.T) {

	service := NewUserService(repository.NewMemoryRepository())

	currentTS := time.Now().UnixNano()
	time.Sleep(time.Nanosecond * 100)

	user, err := service.CreateUser(test_user.User{
		Name: "My new user",
	})
	if err != nil {
		t.Error(err)
		return
	}

	user.Name = "My new name"
	userTmp, _ := service.UpdateUser(user)

	if userTmp.Name != user.Name {
		t.Error("mismatch name")
		return
	}

	if user.UpdatedAt.UnixNano() < currentTS {
		t.Error("Time not set")
		return
	}
}

func TestUserService_DeleteUser(t *testing.T) {

	service := NewUserService(repository.NewMemoryRepository())

	user, err := service.CreateUser(test_user.User{
		Name: "My new user",
	})
	if err != nil {
		t.Error(err)
		return
	}

	deletedUser, err := service.DeleteUser(user.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if deletedUser.ID != user.ID {
		t.Error("Deleted user ID mismatch")
	}
}
