package repository

import (
	"github.com/RaphaCosil/messaging-api/internal/model"
	"gorm.io/gorm"
)

type UserChatRepository interface {
	FindAll() ([]model.UserChat, error)
	Create(userChat model.UserChat) (model.UserChat, error)
	Delete(userChat model.UserChat) error
	FindByUserID(userID uint) ([]model.UserChat, error)
	FindByChatID(chatID uint) ([]model.UserChat, error)
	FindByUserIDAndChatID(userID, chatID uint) (model.UserChat, error)
}

type userChatRepository struct {
	db *gorm.DB
}

func NewUserChatRepository(db *gorm.DB) UserChatRepository {
	return &userChatRepository{db}
}

func (r *userChatRepository) FindAll() ([]model.UserChat, error) {
	var userChats []model.UserChat
	result := r.db.Find(&userChats)
	if result.Error != nil {
		return nil, result.Error
	}
	return userChats, nil
}

func (r *userChatRepository) Create(userChat model.UserChat) (model.UserChat, error) {
	if err := r.db.Create(&userChat).Error; err != nil {
		return model.UserChat{}, err
	}
	return userChat, nil
}

func (r *userChatRepository) Delete(userChat model.UserChat) error {
	if err := r.db.Delete(&userChat).Error; err != nil {
		return err
	}
	return nil
}

func (r *userChatRepository) FindByUserID(userID uint) ([]model.UserChat, error) {
	var userChats []model.UserChat
	if err := r.db.Where("user_id = ?", userID).Find(&userChats).Error; err != nil {
		return nil, err
	}
	return userChats, nil
}

func (r *userChatRepository) FindByChatID(chatID uint) ([]model.UserChat, error) {
	var userChats []model.UserChat
	if err := r.db.Where("chat_id = ?", chatID).Find(&userChats).Error; err != nil {
		return nil, err
	}
	return userChats, nil
}

func (r *userChatRepository) FindByUserIDAndChatID(userID, chatID uint) (model.UserChat, error) {
	var userChat model.UserChat
	if err := r.db.Where("user_id = ? AND chat_id = ?", userID, chatID).First(&userChat).Error; err != nil {
		return model.UserChat{}, err
	}
	return userChat, nil
}
