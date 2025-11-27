package telegrambot

import "sync"

type BotStorageRepository interface {
	GetLastSuccessfulUpdateID() (int, error)
	SaveLastSuccessfulUpdateID(updateID int) error
	SaveMessage(chatID int64, message Message) error
	GetMessagesByChatID(chatID int64, limit int) ([]Message, error)
}

type InmemoryBotStorageRepository struct {
	lastUpdateID int
	messages     map[int64][]Message
	mu           sync.RWMutex
}

func NewInmemoryBotStorageRepository() *InmemoryBotStorageRepository {
	return &InmemoryBotStorageRepository{
		lastUpdateID: -1,
		messages:     make(map[int64][]Message),
		mu:           sync.RWMutex{},
	}
}

func (r *InmemoryBotStorageRepository) GetLastSuccessfulUpdateID() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdateID, nil
}

func (r *InmemoryBotStorageRepository) SaveLastSuccessfulUpdateID(updateID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastUpdateID = updateID
	return nil
}

func (r *InmemoryBotStorageRepository) SaveMessage(chatID int64, message Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages[chatID] = append(r.messages[chatID], message)
	return nil
}

func (r *InmemoryBotStorageRepository) GetMessagesByChatID(chatID int64, limit int) ([]Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	messages, exists := r.messages[chatID]
	if !exists {
		return nil, nil
	}

	if limit > len(messages) {
		limit = len(messages)
	}
	return messages[len(messages)-limit:], nil
}
