package services

import (
	"fmt"
	"homework_2/models"
	"sync"
	"sync/atomic"
)

type UserService struct {
	storage map[int]models.User
	mu      sync.RWMutex
	nextID  int32
}

// NewUserService создает новый экземпляр UserService.
func NewUserService() *UserService {
	return &UserService{
		storage: make(map[int]models.User),
		nextID:  1, // Начинаем ID с 1
	}
}

// Create сохраняет нового пользователя и генерирует ID.
func (s *UserService) Create(user models.User) models.User {
	fmt.Print("Зарегестрированно обращение")
	s.mu.Lock()
	defer s.mu.Unlock()

	// Атомарно увеличиваем и получаем новый ID
	newID := atomic.AddInt32(&s.nextID, 1)
	user.ID = int(newID - 1) // ID перед инкрементом

	s.storage[user.ID] = user
	return user
}

// GetAll возвращает список всех пользователей.
func (s *UserService) GetAll() []models.User {
	s.mu.RLock() // RLock для чтения
	defer s.mu.RUnlock()

	users := make([]models.User, 0, len(s.storage))
	for _, user := range s.storage {
		users = append(users, user)
	}
	return users
}

// GetByID возвращает пользователя по ID.
func (s *UserService) GetByID(id int) (models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.storage[id]
	return user, ok
}

// Update обновляет существующего пользователя.
func (s *UserService) Update(id int, updatedUser models.User) (models.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.storage[id]; !ok {
		return models.User{}, false // Пользователь не найден
	}

	updatedUser.ID = id // Гарантируем, что ID не изменится
	s.storage[id] = updatedUser
	return updatedUser, true
}

// Delete удаляет пользователя по ID.
func (s *UserService) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.storage[id]; !ok {
		return false // Пользователь не найден
	}

	delete(s.storage, id)
	return true
}
