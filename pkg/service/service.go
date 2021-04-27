package service

import (
	test_user "test-user"
	"test-user/pkg/repository"
	"time"
)

type UserNameEmptyError struct {
}

func (err *UserNameEmptyError) Error() string {
	return "user name empty"
}

type UserService interface {
	CreateUser(user test_user.User) (test_user.User, error)
	GetUser(id string) (test_user.User, error)
	GetUsers() ([]test_user.User, error)
	DeleteUser(id string) (test_user.User, error)
	UpdateUser(user test_user.User) (test_user.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) CreateUser(user test_user.User) (test_user.User, error) {

	if user.Name == "" {
		return user, &UserNameEmptyError{}
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.repo.Create(user)
}

func (s *userService) GetUser(id string) (test_user.User, error) {
	return s.repo.Get(id)
}

func (s *userService) GetUsers() ([]test_user.User, error) {
	return s.repo.GetAll()
}

func (s *userService) DeleteUser(id string) (test_user.User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return user, err
	}
	return user, s.repo.Delete(id)
}

func (s *userService) UpdateUser(user test_user.User) (test_user.User, error) {
	userTmp, err := s.repo.Get(user.ID)
	if err != nil {
		return user, err
	}
	user.UpdatedAt = time.Now()
	user.CreatedAt = userTmp.CreatedAt
	return user, s.repo.Update(user)
}