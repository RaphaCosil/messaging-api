package model

import (
	"time"
)

type CustomerChat struct {
    CustomerID   uint      `gorm:"primaryKey;column:customer_id"`
    ChatID   uint      `gorm:"primaryKey;column:chat_id"`
    Role     string    `gorm:"default:participant;not null"`
    JoinedAt time.Time `gorm:"column:joined_at;autoCreateTime"`

    Customer     Customer      `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Chat     Chat      `gorm:"foreignKey:ChatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}