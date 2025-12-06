package app

import (
	"log/slog"
	"speaking-club-bot/internal/client/ai"
	"speaking-club-bot/internal/client/telegram"
	"speaking-club-bot/internal/config"
	"speaking-club-bot/internal/storage"
	"speaking-club-bot/internal/bot"
)

type App struct {
	bot    bot.Bot
	config *config.Config
}

func (a *App) App() bot.Bot {
	return a.bot
}

func (a *App) Confgi() *config.Config {
	return a.config
}

func New(config *config.Config) (App, error) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	storageRepository := storage.NewInmemory()
	mentorClient := ai.New(config.AIProviderAPIBaseURL, config.AIProviderToken, config.MentorSystemPrompt)
	telegramClient := telegram.New(config.TelegramAPIBaseURL, config.TelegramBotToken)

	bot := bot.NewBot(telegramClient, storageRepository, mentorClient)

	return App{
		bot:    bot,
		config: config,
	}, nil
}
