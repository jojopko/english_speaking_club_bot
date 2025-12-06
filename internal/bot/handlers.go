package bot

import (
	"log/slog"
	"speaking-club-bot/internal/client/telegram"
	"speaking-club-bot/internal/config"
	"strings"
)

func ProcessUpdate(
	bot *Bot,
	config *config.Config,
	update *telegram.Update,
) error {

	chatId, isAuth := isAuth(update, config.WhitelistIDs)
	if !isAuth {
		slog.Debug("Ignoring message from unauthorized chat ID", "id", chatId)
		return nil
	}

	var err error
	if update.Message != nil {
		switch {
		case (strings.HasPrefix(update.Message.Text, "/helpme") || strings.HasPrefix(update.Message.Text, config.TelegramBotNickname)) && update.Message.ReplyTo != nil:
			err = helpmeHandle(bot, update)
		case strings.HasPrefix(update.Message.Text, config.TelegramBotNickname) && update.Message.ReplyTo == nil:
			err = askHandler(bot, update)
		}
	}
	return err
}

func isAuth(update *telegram.Update, whitelist map[int64]struct{}) (int64, bool) {
	if update.Message != nil {
		_, ok := whitelist[update.Message.Chat.ID]
		return update.Message.Chat.ID, ok
	}
	return 0, false
}
