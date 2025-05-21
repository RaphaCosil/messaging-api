package service

import(
	"github.com/RaphaCosil/messaging-api/internal/model"
	"github.com/RaphaCosil/messaging-api/internal/repository"
)

type UserChatService interface {
	FindAll() ([]model.UserChat, error)
	Create(userID, chatID uint) (model.UserChat, error)
	Delete(userID, chatID uint) error
	FindByUserID(userID uint) ([]model.UserChat, error)
	FindByChatID(chatID uint) ([]model.UserChat, error)
	UserHasAccessToChat(userID, chatID uint) (bool)
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

func (s *userChatService) Create(userID, chatID uint) (model.UserChat, error) {
	userChat := model.UserChat{
		UserID: userID,
		ChatID: chatID,
	}
	return s.userChatRepo.Create(userChat)
}

func (s *userChatService) Delete(userID, chatID uint) error {
	userChat := model.UserChat{
		UserID: userID,
		ChatID: chatID,
	}
	return s.userChatRepo.Delete(userChat)
}

func (s *userChatService) FindByUserID(userID uint) ([]model.UserChat, error) {
	return s.userChatRepo.FindByUserID(userID)
}

func (s *userChatService) FindByChatID(chatID uint) ([]model.UserChat, error) {
	return s.userChatRepo.FindByChatID(chatID)
}

func (s *userChatService) UserHasAccessToChat(userID, chatID uint) (bool) {
	_, err := s.userChatRepo.FindByUserIDAndChatID(userID, chatID)
	if err != nil {
		return false
	}
	return true
}
