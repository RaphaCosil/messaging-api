package model

import (
	"time"
)

type Customer struct {
    CustomerID    uint      `gorm:"primaryKey;column:customer_id"`
    Username  string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}