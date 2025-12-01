package app

import (
	"context"
	"log/slog"
	"speaking-club-bot/internal/client/telegram"
	"speaking-club-bot/internal/config"
	"speaking-club-bot/internal/telegrambot"
)

func (a *App) Start(ctx context.Context) error {
	bot := &a.bot

	offset, err := bot.StorageRepository().GetLastSuccessfulUpdateID()
	if err != nil {
		return err
	}

	errorChan := make(chan error)
	defer close(errorChan)

	go errorHandler(errorChan, bot, a.config.AdminID)

	slog.Info("Starting bot with offset:", "offset", offset)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down bot...")
			return nil
		default:
			updates, err := bot.TelegramClient().GetUpdates(telegram.GetUpdatesRequest{
				Offset: offset,
				Limit:  60,
			})
			if err != nil {
				slog.Error("Error getting updates:", "err", err)
				continue
			}

			for _, update := range updates {
				slog.Debug("Received update:", "update", update)
				offset = update.UpdateId + 1
				if update.Message == nil {
					continue
				}

				go processUpdate(&a.bot, a.config, &update, errorChan)
			}
		}
	}
}

func errorHandler(errorChan chan error, bot *telegrambot.Bot, adminId int64) {
	for err := range errorChan {
		errorTextMessage := "Error in bot loop: " + err.Error()
		slog.Error("Error received:", "err", err)
		_, err := bot.TelegramClient().SendMessage(telegram.SendMessageRequest{
			ChatID: adminId,
			Text:   errorTextMessage,
		})
		if err != nil {
			slog.Warn("Failed to send error message to admin:", "err", err)
		}
	}
}

func processUpdate(
	bot *telegrambot.Bot,
	config *config.Config,
	update *telegram.Update,
	errorChan chan error,
) error {
	err := telegrambot.ProcessUpdate(bot, config, update)
	if err != nil {
		errorChan <- err
		return nil
	}
	if err := bot.StorageRepository().SaveLastSuccessfulUpdateID(update.UpdateId); err != nil {
		errorChan <- err
	}
	return err
}
