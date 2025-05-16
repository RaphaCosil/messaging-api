package model

import (
    "time"
    "gorm.io/gorm"
)

type Chat struct {
    ChatID    uint      `gorm:"primaryKey;column:chat_id"`
    ChatName  string    `gorm:"not null"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

    Users     []User    `gorm:"many2many:user_chats"`
    Messages  []Message `gorm:"foreignKey:ChatID"`
}