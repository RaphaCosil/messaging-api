package service

import (
	"github.com/RaphaCosil/messaging-api/internal/model"
	"github.com/RaphaCosil/messaging-api/internal/repository"
)

type UserService interface {
	FindAll() ([]model.Customer, error)
	FindByID(id uint) (model.Customer, error)
	Create(user model.Customer) (model.Customer, error)
	Update(id uint, user model.Customer) (model.Customer, error)
	Delete(id uint) error
	FindByUsername(username string) (model.Customer, error)
	FindByChatID(chatID uint) ([]model.Customer, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) FindAll() ([]model.Customer, error) {
	return s.userRepo.FindAll()
}

func (s *userService) FindByID(id uint) (model.Customer, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) Create(user model.Customer) (model.Customer, error) {
	return s.userRepo.Create(user)
}

func (s *userService) Update(id uint, user model.Customer) (model.Customer, error) {
	return s.userRepo.Update(id, user)
}

func (s *userService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) FindByUsername(username string) (model.Customer, error) {
	return s.userRepo.FindByUsername(username)
}

func (s *userService) FindByChatID(chatID uint) ([]model.Customer, error) {
	return s.userRepo.FindByChatID(chatID)
}
