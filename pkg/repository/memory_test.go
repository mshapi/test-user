package repository

import (
	"fmt"
	"sync"
	test_user "test-user"
	"testing"
	"time"
)

func generateUsers(count int) []test_user.User {
	var users = make([]test_user.User, count, count) //nolint:gosimple
	for i := 0; i < count; i++ {
		users[i] = test_user.User{
			Name:      fmt.Sprintf("Test User %s %d", randString(5), i),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}
	}
	return users
}

const countTestUsers = 200

func TestMemoryRepository_Create(t *testing.T) {

	users := generateUsers(countTestUsers)

	wg := &sync.WaitGroup{}
	wg.Add(len(users))

	repo := NewMemoryRepository()

	for i := range users {
		go func(n int) {
			users[n], _ = repo.Create(users[n])
			wg.Done()
		}(i)
	}

	wg.Wait()

	t.Run("Check creating users", func(t *testing.T) {
		for _, user := range users {
			userRepo, err := repo.Get(user.ID)
			if err != nil {
				t.Error(err)
				continue
			}
			if userRepo.Name != user.Name {
				t.Error("Username mismatch")
			}
		}
	})
}

func TestMemoryRepository_Get(t *testing.T) {

	users := generateUsers(countTestUsers)[:]

	tests := []struct {
		name  string
		users []test_user.User
		repo  UserRepository
		check func(repo UserRepository, users ...test_user.User) bool
	}{
		{
			name:  "OK",
			users: nil,
			repo: func() UserRepository {
				repo := NewMemoryRepository()
				wg := &sync.WaitGroup{}
				wg.Add(len(users))
				for _, user := range users {
					go func(user test_user.User) {
						_, _ = repo.Create(user)
						wg.Done()
					}(user)
				}
				wg.Wait()
				return repo
			}(),
			check: func(repo UserRepository, users ...test_user.User) bool {
				for _, user := range users {
					_, err := repo.Get(user.ID)
					if err != nil {
						return false
					}
				}
				return true
			},
		},
		{
			name:  "FAIL",
			users: users,
			repo:  NewMemoryRepository(),
			check: func(repo UserRepository, users ...test_user.User) bool {
				if len(users) <= 0 {
					return true
				}
				for _, user := range users {
					_, err := repo.Get(user.ID)
					if err == nil {
						return false
					}
					_, ok := err.(*UserNotFoundError)
					return ok //nolint:staticcheck
				}
				return true
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var tmpUsers = testCase.users
			if len(tmpUsers) <= 0 {
				tmpUsers, _ = testCase.repo.GetAll()
			}
			if !testCase.check(testCase.repo, tmpUsers...) {
				t.Error("error")
			}
		})
	}
}

func TestMemoryRepository_Delete(t *testing.T) {
	user := generateUsers(1)[0]

	repo := NewMemoryRepository()

	user, _ = repo.Create(user)

	_, err := repo.Get(user.ID)

	if err != nil {
		t.Error(err)
	}

	err = repo.Delete(user.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestMemoryRepository_Update(t *testing.T) {
	user := generateUsers(1)[0]

	repo := NewMemoryRepository()

	user, _ = repo.Create(user)

	_, err := repo.Get(user.ID)

	if err != nil {
		t.Error(err)
	}

	user.Name = "New name"

	err = repo.Update(user)
	if err != nil {
		t.Error(err)
	}

	updatedUser, _ := repo.Get(user.ID)
	if updatedUser.Name != user.Name {
		t.Error("user not updated")
	}
}
