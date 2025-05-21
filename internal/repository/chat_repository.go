package repository

import (
	"github.com/RaphaCosil/messaging-api/internal/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	FindAll() ([]model.Chat, error)
	FindByID(id uint) (model.Chat, error)
	Create(chat model.Chat) (model.Chat, error)
	Update(id uint, chat model.Chat) (model.Chat, error)
	Delete(id uint) error
	FindByUserID(userID uint) ([]model.Chat, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db}
}

func (r *chatRepository) FindAll() ([]model.Chat, error) {
	var chats []model.Chat
	result := r.db.Find(&chats)
	if result.Error != nil {
		return nil, result.Error
	}
	return chats, nil
}

func (r *chatRepository) FindByID(id uint) (model.Chat, error) {
	var chat model.Chat
	if err := r.db.First(&chat, id).Error; err != nil {
		return model.Chat{}, err
	}
	return chat, nil
}

func (r *chatRepository) Create(chat model.Chat) (model.Chat, error) {
	if err := r.db.Create(&chat).Error; err != nil {
		return model.Chat{}, err
	}
	return chat, nil
}

func (r *chatRepository) Update(id uint, chat model.Chat) (model.Chat, error) {
	if err := r.db.Model(&chat).Where("chat_id = ?", id).Updates(chat).Error; err != nil {
		return model.Chat{}, err
	}
	return chat, nil
}

func (r *chatRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Chat{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) FindByUserID(userID uint) ([]model.Chat, error) {
	var chats []model.Chat
	if err := r.db.Table("user_chats").Select("chats.*").
		Joins("JOIN chats ON user_chats.chat_id = chats.chat_id").
		Where("user_chats.user_id = ?", userID).
		Scan(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}
