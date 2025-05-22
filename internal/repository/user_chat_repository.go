package repository

import (
	"github.com/RaphaCosil/messaging-api/internal/model"
	"gorm.io/gorm"
)

type UserChatRepository interface {
	FindAll() ([]model.CustomerChat, error)
	Create(userChat model.CustomerChat) (model.CustomerChat, error)
	Delete(userChat model.CustomerChat) error
	FindByUserID(userID uint) ([]model.CustomerChat, error)
	FindByChatID(chatID uint) ([]model.CustomerChat, error)
	FindByUserIDAndChatID(userID, chatID uint) (model.CustomerChat, error)
}

type userChatRepository struct {
	db *gorm.DB
}

func NewUserChatRepository(db *gorm.DB) UserChatRepository {
	return &userChatRepository{db}
}

func (r *userChatRepository) FindAll() ([]model.CustomerChat, error) {
	var userChats []model.CustomerChat
	result := r.db.Find(&userChats)
	if result.Error != nil {
		return nil, result.Error
	}
	return userChats, nil
}

func (r *userChatRepository) Create(userChat model.CustomerChat) (model.CustomerChat, error) {
	if err := r.db.Create(&userChat).Error; err != nil {
		return model.CustomerChat{}, err
	}
	return userChat, nil
}

func (r *userChatRepository) Delete(userChat model.CustomerChat) error {
	if err := r.db.Delete(&userChat).Error; err != nil {
		return err
	}
	return nil
}

func (r *userChatRepository) FindByUserID(userID uint) ([]model.CustomerChat, error) {
	var userChats []model.CustomerChat
	if err := r.db.Where("customer_id = ?", userID).Find(&userChats).Error; err != nil {
		return nil, err
	}
	return userChats, nil
}

func (r *userChatRepository) FindByChatID(chatID uint) ([]model.CustomerChat, error) {
	var userChats []model.CustomerChat
	if err := r.db.Where("chat_id = ?", chatID).Find(&userChats).Error; err != nil {
		return nil, err
	}
	return userChats, nil
}

func (r *userChatRepository) FindByUserIDAndChatID(userID, chatID uint) (model.CustomerChat, error) {
	var userChat model.CustomerChat
	if err := r.db.Where("customer_id = ? AND chat_id = ?", userID, chatID).First(&userChat).Error; err != nil {
		return model.CustomerChat{}, err
	}
	return userChat, nil
}
