package bot

import (
	"speaking-club-bot/internal/client/ai"
	"speaking-club-bot/internal/client/telegram"
	"speaking-club-bot/internal/storage"
)

type Bot struct {
	telegramClient    telegram.Client
	storageRepository storage.Repository
	mentorRepository  ai.Client
}

func NewBot(
	telegramClient telegram.Client,
	storageRepo storage.Repository,
	mentorRepo ai.Client,
) Bot {
	return Bot{
		telegramClient:    telegramClient,
		storageRepository: storageRepo,
		mentorRepository:  mentorRepo,
	}
}

func (b *Bot) StorageRepository() storage.Repository {
	return b.storageRepository
}

func (b *Bot) MentorClient() ai.Client {
	return b.mentorRepository
}

func (b *Bot) TelegramClient() telegram.Client {
	return b.telegramClient
}
