package repository

import (
	"github.com/RaphaCosil/messaging-api/internal/model"
	"gorm.io/gorm"
)

type MessageRepository interface {
	FindAll() ([]model.Message, error)
	FindByID(id uint) (model.Message, error)
	Create(message model.Message) (model.Message, error)
	Delete(id uint) error
	FindByChatID(chatID uint) ([]model.Message, error)
	FindByUserID(userID uint) ([]model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) FindAll() ([]model.Message, error) {
	var messages []model.Message
	result := r.db.Preload("Chat").Preload("User").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (r *messageRepository) FindByID(id uint) (model.Message, error) {
	var message model.Message
	if err := r.db.Preload("Chat").Preload("User").First(&message, id).Error; err != nil {
		return model.Message{}, err
	}
	return message, nil
}

func (r *messageRepository) Create(message model.Message) (model.Message, error) {
	if err := r.db.Create(&message).Error; err != nil {
		return model.Message{}, err
	}
	return message, nil
}

func (r *messageRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Message{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *messageRepository) FindByChatID(chatID uint) ([]model.Message, error) {
	var messages []model.Message
	result := r.db.Where("chat_id = ?", chatID).Preload("Chat").Preload("User").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (r *messageRepository) FindByUserID(userID uint) ([]model.Message, error) {
	var messages []model.Message
	result := r.db.Where("user_id = ?", userID).Preload("Chat").Preload("User").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
