package telegrambot

import (
	"fmt"
	"log"
	"speaking-club-bot/internal/client/telegram"
	"speaking-club-bot/internal/config"
	"strings"
)

type HandlerFunc struct {
	name           string
	CheckCondition func(bot *Bot, message *telegram.Message) bool
	Handle         func(bot *Bot, message *telegram.Message) error
}

func isAuth(whitelist map[int64]struct{}, telegramId int64) bool {
	_, ok := whitelist[telegramId]
	return ok
}

func ProcessUpdate(
	bot *Bot,
	config *config.Config,
	update *telegram.Update,
) error {
	message := update.Message
	isAuth := isAuth(config.WhitelistIDs, message.Chat.ID)
	if !isAuth {
		log.Printf("Ignoring message from unauthorized chat ID: %d", message.Chat.ID)
		return nil
	}

	handlers := []HandlerFunc{
		{
			name: "LoggerHandler",
			CheckCondition: func(bot *Bot, message *telegram.Message) bool {
				return true
			},
			Handle: EchoHandler,
		},
		{
			name: "Help me command handler",
			CheckCondition: func(bot *Bot, message *telegram.Message) bool {
				return strings.HasPrefix(message.Text, "/helpme")
			},
			Handle: HelpMeCommandHandler,
		},
	}

	for _, handler := range handlers {
		if handler.CheckCondition(bot, message) {
			log.Printf("Handler %s condition met, processing...", handler.name)
			if err := handler.Handle(bot, message); err != nil {
				log.Printf("Error in handler %s: %v", handler.name, err)
				return err
			}
		}
	}

	return nil
}

func EchoHandler(bot *Bot, message *telegram.Message) error {
	log.Printf("Echoing message: %s from %s", message.Text, message.From.Name())
	return nil
}

const infoMessage = "Use this command to reply to the message you need help on"

func HelpMeCommandHandler(bot *Bot, message *telegram.Message) error {
	if message.ReplyTo == nil {
		_, err := bot.TelegramClient().SendMessage(telegram.SendMessageRequest{
			ChatID: message.Chat.ID,
			Text:   infoMessage,
		})
		return err
	}

	commandSlice := strings.SplitN(message.Text, " ", 2)
	var extraQuestion string
	if len(commandSlice) < 2 {
		extraQuestion = ""
	} else {
		extraQuestion = commandSlice[1]
	}
	extraQuestion = strings.TrimSpace(extraQuestion)

	s := message.ReplyTo.Text
	answer, err := bot.MentorClient().AskQuestion(s, extraQuestion)
	if err != nil {
		return fmt.Errorf("MentorRepo.AskToFixMistake: %w", err)
	}

	_, err = bot.TelegramClient().SendMessage(telegram.SendMessageRequest{
		ChatID: message.Chat.ID,
		Text:   answer,
		ReplyParameters: &telegram.ReplyParameters{
			MessageID: message.MessageID,
		},
	})
	return err
}
