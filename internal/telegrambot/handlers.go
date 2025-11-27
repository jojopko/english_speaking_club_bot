package telegrambot

import (
	"log"
	"strings"
)

type HandlerFunc struct {
	name           string
	CheckCondition func(bot *Bot, message *Message) bool
	Handle         func(bot *Bot, message *Message) error
}

func BaseMiddleware(bot *Bot, message *Message) error {
	isAuth := false
	for _, id := range bot.Whitelist {
		if message.Chat.ID == id {
			isAuth = true
			break
		}
	}
	if !isAuth {
		log.Printf("Ignoring message from unauthorized chat ID: %d", message.Chat.ID)
		return nil
	}

	handlers := []HandlerFunc{
		{
			name: "LoggerHandler",
			CheckCondition: func(bot *Bot, message *Message) bool {
				return true
			},
			Handle: EchoHandler,
		},
		{
			name: "Help me command handler",
			CheckCondition: func(bot *Bot, message *Message) bool {
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

func EchoHandler(bot *Bot, message *Message) error {
	log.Printf("Echoing message: %s from %s", message.Text, message.From.Name())
	return nil
}

func MessageSaverHandler(bot *Bot, message *Message) error {
	s := strings.TrimSpace(message.Text)
	if s == "" {
		return nil
	}

	bot.StorageRepo.SaveMessage(message.Chat.ID, *message)
	return nil
}

func HelpMeCommandHandler(bot *Bot, message *Message) error {
	extraQuestion := ""
	if message.ReplyTo != nil {
		s := message.ReplyTo.Text
		answer, err := bot.MentorRepo.AskToFixMistake(s, extraQuestion)
		if err != nil {
			return err
		}
		bot.SendMessage(message.Chat.ID, answer, message.MessageID)
	} else {
		const infoMessage = "Use this command to reply to the message you need help on"
		bot.SendMessage(message.Chat.ID, infoMessage, 0)
	}
	return nil
}
