package service

import(
	"github.com/RaphaCosil/messaging-api/internal/model"
	"github.com/RaphaCosil/messaging-api/internal/repository"
)

type UserChatService interface {
	FindAll() ([]model.UserChat, error)
	Create(userChat model.UserChat) (model.UserChat, error)
	Delete(id uint) error
	FindByUserID(userID uint) ([]model.UserChat, error)
	FindByChatID(chatID uint) ([]model.UserChat, error)
	FindByUserIDAndChatID(userID, chatID uint) (model.UserChat, error)
}

type userChatService struct {
	userChatRepo repository.UserChatRepository
}

func NewUserChatService(userChatRepo repository.UserChatRepository) UserChatService {
	return &userChatService{userChatRepo}
}

func (s *userChatService) FindAll() ([]model.UserChat, error) {
	return s.userChatRepo.FindAll()
}

func (s *userChatService) Create(userChat model.UserChat) (model.UserChat, error) {
	return s.userChatRepo.Create(userChat)
}

func (s *userChatService) Delete(id uint) error {
	return s.userChatRepo.Delete(id)
}

func (s *userChatService) FindByUserID(userID uint) ([]model.UserChat, error) {
	return s.userChatRepo.FindByUserID(userID)
}

func (s *userChatService) FindByChatID(chatID uint) ([]model.UserChat, error) {
	return s.userChatRepo.FindByChatID(chatID)
}

func (s *userChatService) FindByUserIDAndChatID(userID, chatID uint) (model.UserChat, error) {
	return s.userChatRepo.FindByUserIDAndChatID(userID, chatID)
}
