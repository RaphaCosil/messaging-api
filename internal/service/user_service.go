package service

import(
	"github.com/RaphaCosil/messaging-api/internal/model"
	"github.com/RaphaCosil/messaging-api/internal/repository"
)

type UserService interface {
	FindAll() ([]model.User, error)
	FindByID(id uint) (model.User, error)
	Create(user model.User) (model.User, error)
	Update(user model.User) (model.User, error)
	Delete(id uint) error
	FindByUsername(username string) (model.User, error)
	FindByChatID(chatID uint) ([]model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) FindAll() ([]model.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) FindByID(id uint) (model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) Create(user model.User) (model.User, error) {
	return s.userRepo.Create(user)
}

func (s *userService) Update(user model.User) (model.User, error) {
	return s.userRepo.Update(user)
}

func (s *userService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *userService) FindByUsername(username string) (model.User, error) {
	return s.userRepo.FindByUsername(username)
}

func (s *userService) FindByChatID(chatID uint) ([]model.User, error) {
	return s.userRepo.FindByChatID(chatID)
}
