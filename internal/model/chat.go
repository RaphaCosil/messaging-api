package model

import (
    "time"
)

type Chat struct {
    ChatID    uint      `gorm:"primaryKey;column:chat_id"`
    ChatName  string    `gorm:"not null"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}