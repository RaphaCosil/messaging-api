package model

import (
    "time"
)

type Message struct {
    MessageID uint      `gorm:"primaryKey;column:message_id"`
    ChatID    uint      `gorm:"not null"`
    UserID    uint      `gorm:"not null"`
    Content   string    `gorm:"type:text;not null"`
    SentAt    time.Time `gorm:"column:sent_at;autoCreateTime"`

    Chat      Chat      `gorm:"foreignKey:ChatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}