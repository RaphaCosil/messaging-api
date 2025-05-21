package model

import "time"

type GenericMessage struct {
    Type      string      `json:"type"`
    Content   interface{} `json:"content"`
    Timestamp time.Time   `json:"timestamp"`
}