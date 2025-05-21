package model

import (
	"sync"
)

type Hub struct {
    clients   map[*Client]bool
    broadcast chan GenericMessage
    mu        sync.Mutex
}