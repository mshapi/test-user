package repository

import (
	"sync"
	"test-user"
)

const idLength uint8 = 8

type memoryRepository struct {
	users map[string]test_user.User
	mutex *sync.RWMutex
}

func NewMemoryRepository() UserRepository {
	return &memoryRepository{
		users: make(map[string]test_user.User),
		mutex: new(sync.RWMutex),
	}
}

func (m *memoryRepository) Create(user test_user.User) (test_user.User, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user.ID = randString(idLength)

	m.users[user.ID] = user

	return user, nil
}

func (m *memoryRepository) Get(id string) (test_user.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var user test_user.User
	var ok bool

	user, ok = m.users[id]
	if !ok {
		return user, &UserNotFoundError{id}
	}
	return user, nil
}

func (m *memoryRepository) GetAll() ([]test_user.User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	users := make([]test_user.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *memoryRepository) Delete(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.users, id)
	return nil
}

func (m *memoryRepository) Update(user test_user.User) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.users[user.ID]; ok {
		m.users[user.ID] = user
	}

	return nil
}
