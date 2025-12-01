package storage

import (
	"speaking-club-bot/internal/client/telegram"
	"sync"
)

type Repository interface {
	GetLastSuccessfulUpdateID() (int64, error)
	SaveLastSuccessfulUpdateID(updateID int64) error
}

type repository struct {
	lastUpdateID int64
	messages     map[int64][]telegram.Message
	mu           *sync.RWMutex
}

func NewInmemory() Repository {
	return &repository{
		lastUpdateID: -1,
		messages:     make(map[int64][]telegram.Message),
		mu:           &sync.RWMutex{},
	}
}

func (r *repository) GetLastSuccessfulUpdateID() (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdateID, nil
}

func (r *repository) SaveLastSuccessfulUpdateID(updateID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastUpdateID = updateID
	return nil
}
