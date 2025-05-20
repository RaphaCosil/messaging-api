package service

import(
	"github.com/RaphaCosil/messaging-api/internal/model"
	"github.com/RaphaCosil/messaging-api/internal/repository"
)

type ChatService interface {
	FindAll() ([]model.Chat, error)
	FindByID(id uint) (model.Chat, error)
	Create(chat model.Chat) (model.Chat, error)
	Update(chat model.Chat) (model.Chat, error)
	Delete(id uint) error
	FindByUserID(userID uint) ([]model.Chat, error)
}

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{chatRepo}
}

func (s *chatService) FindAll() ([]model.Chat, error) {
	return s.chatRepo.FindAll()
}

func (s *chatService) FindByID(id uint) (model.Chat, error){
	return s.chatRepo.FindByID(id)
}

func (s *chatService) Create(chat model.Chat) (model.Chat, error){
	return s.chatRepo.Create(chat)
}

func (s *chatService) Update(chat model.Chat) (model.Chat, error){
	return s.chatRepo.Update(chat)
}

func (s *chatService) Delete(id uint) error {
	return s.chatRepo.Delete(id)
}

func (s *chatService) FindByUserID(userID uint) ([]model.Chat, error){
	return s.chatRepo.FindByUserID(userID)
}
