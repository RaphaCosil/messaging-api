package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    UserID    uint      `gorm:"primaryKey;column:user_id"`
    Username  string    `gorm:"unique;not null"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

    Chats     []Chat    `gorm:"many2many:user_chats"`
    Messages  []Message `gorm:"foreignKey:UserID"`
}