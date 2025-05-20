package model

import (
    "time"
)

type User struct {
    UserID    uint      `gorm:"primaryKey;column:user_id"`
    Username  string    `gorm:"unique;not null"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}