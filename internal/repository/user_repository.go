package repository

import (
	"github.com/RaphaCosil/messaging-api/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.Customer, error)
	FindByID(id uint) (model.Customer, error)
	Create(user model.Customer) (model.Customer, error)
	Update(id uint, user model.Customer) (model.Customer, error)
	Delete(id uint) error
	FindByUsername(username string) (model.Customer, error)
	FindByChatID(chatID uint) ([]model.Customer, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Login(username, password string) (model.Customer, error) {
	var user model.Customer
	if err := r.db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		return model.Customer{}, err
	}
	return user, nil
}

func (r *userRepository) FindAll() ([]model.Customer, error) {
	var users []model.Customer
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) FindByID(id uint) (model.Customer, error) {
	var user model.Customer
	if err := r.db.First(&user, id).Error; err != nil {
		return model.Customer{}, err
	}
	return user, nil
}

func (r *userRepository) Create(user model.Customer) (model.Customer, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return model.Customer{}, err
	}
	return user, nil
}

func (r *userRepository) Update(id uint, user model.Customer) (model.Customer, error) {
	if err := r.db.Model(&user).Where("customer_id = ?", id).Updates(user).Error; err != nil {
		return model.Customer{}, err
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Customer{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByUsername(username string) (model.Customer, error) {
	var user model.Customer
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return model.Customer{}, err
	}
	return user, nil
}

func (r *userRepository) FindByChatID(chatID uint) ([]model.Customer, error) {
	var users []model.Customer
	if err := r.db.Table("customer_chats").Select("customers.*").
		Joins("JOIN customers ON customer_chats.customer_id = customers.customer_id").
		Where("customer_chats.chat_id = ?", chatID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
