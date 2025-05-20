package service

import(
	"github.com/RaphaCosil/messaging-api/internal/model"
	"github.com/RaphaCosil/messaging-api/internal/repository"
)

type MessageService interface {
	FindAll() ([]model.Message, error)
	FindByID(id uint) (model.Message, error)
	Create(message model.Message) (model.Message, error)
	Update(message model.Message) (model.Message, error)
	Delete(id uint) error
	FindByChatID(chatID uint) ([]model.Message, error)
	FindByUserID(userID uint) ([]model.Message, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService(messageRepo repository.MessageRepository) MessageService {
	return &messageService{messageRepo}
}

func (s *messageService) FindAll() ([]model.Message, error) {
	return s.messageRepo.FindAll()
}

func (s *messageService) FindByID(id uint) (model.Message, error){
	return s.messageRepo.FindByID(id)
}

func (s *messageService) Create(message model.Message) (model.Message, error){
	return s.messageRepo.Create(message)
}

func (s *messageService) Update(message model.Message) (model.Message, error){
	return s.messageRepo.Update(message)
}

func (s *messageService) Delete(id uint) error {
	return s.messageRepo.Delete(id)
}

func (s *messageService) FindByChatID(chatID uint) ([]model.Message, error){
	return s.messageRepo.FindByChatID(chatID)
}

func (s *messageService) FindByUserID(userID uint) ([]model.Message, error){
	return s.messageRepo.FindByUserID(userID)
}
