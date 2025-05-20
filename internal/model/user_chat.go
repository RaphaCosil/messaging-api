package model

import (
    "time"
)

type UserChat struct {
    UserID   uint      `gorm:"primaryKey;column:user_id"`
    ChatID   uint      `gorm:"primaryKey;column:chat_id"`
    Role     string    `gorm:"default:participant;not null"`
    JoinedAt time.Time `gorm:"column:joined_at;autoCreateTime"`

    User     User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Chat     Chat      `gorm:"foreignKey:ChatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}