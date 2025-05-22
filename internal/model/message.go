package model

import (
	"time"
)

type Message struct {
    MessageID uint      `gorm:"primaryKey;column:message_id"`
    ChatID    uint      `gorm:"not null"`
    CustomerID    uint      `gorm:"not null"`
    Content   string    `gorm:"type:text;not null"`
    SentAt    time.Time `gorm:"column:sent_at;autoCreateTime"`
}