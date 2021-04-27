package repository

import (
	"fmt"
	"math/rand"
	test_user "test-user"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type UserRepository interface {
	Create(user test_user.User) (test_user.User, error)
	Get(id string) (test_user.User, error)
	GetAll() ([]test_user.User, error)
	Delete(id string) error
	Update(user test_user.User) error
}

const bytesForRandString = "abcdefghijklmnopqrstuvwxyz0123456789"

func randString(length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = bytesForRandString[rand.Intn(len(bytesForRandString))]
	}
	return string(b)
}

type UserNotFoundError struct {
	ID string
}

func (err *UserNotFoundError) Error() string {
	return fmt.Sprintf("User '%s' not found", err.ID)
}
