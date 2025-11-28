package main

import (
	"log"
	"os"
	"speaking-club-bot/internal/ai"
	"speaking-club-bot/internal/telegrambot"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	whitelist := whitelistIdsToSlice(os.Getenv("WHITELIST_IDS"))
	log.Printf("Whitelist: %v", whitelist)
	if len(whitelist) == 0 {
		log.Fatal("No IDs in WHITELIST")
	}

	adminId, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	aiProviderToken := os.Getenv("AI_TOKEN")

	storageRepo := telegrambot.NewInmemoryBotStorageRepository()
	mentorRepo := ai.NewAIMentorRepository(aiProviderToken)
	bot, err := telegrambot.NewBot(telegramBotToken, whitelist, adminId, storageRepo, mentorRepo)
	if err != nil {
		log.Fatal("Failed to create Telegram bot:", err)
	}

	if err := startBot(bot); err != nil {
		log.Fatal("Failed:", err)
	}
}

func whitelistIdsToSlice(whitelistIds string) []int64 {
	whitelistCleaned := strings.TrimSpace(whitelistIds)
	if len(whitelistCleaned) == 0 {
		return nil
	}

	whitelist := make([]int64, 0, 10)
	whitelistStrings := strings.SplitSeq(whitelistCleaned, ",")
	for s := range whitelistStrings {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		} else {
			whitelist = append(whitelist, n)
		}
	}
	return whitelist
}

func startBot(bot *telegrambot.Bot) error {
	offset, err := bot.StorageRepo.GetLastSuccessfulUpdateID()
	if err != nil {
		return err
	}

	errorChan := make(chan error)

	go func() {
		for err := range errorChan {
			log.Printf("Error in bot loop: %v", err)
			errorTextMessage := "Error in bot loop: " + err.Error()
			bot.SendMessage(bot.AdminId, errorTextMessage, 0)
		}
	}()

	log.Printf("Starting bot with offset: %d", offset)

	for {
		updates, err := bot.GetUpdates(offset)
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			continue
		}

		for _, update := range updates {
			log.Printf("Received update: %+v", update)
			offset = update.UpdateID + 1
			go func(update telegrambot.Update) {
				if update.Message == nil {
					return
				}
				err := telegrambot.BaseMiddleware(bot, update.Message)
				if err != nil {
					errorChan <- err
				} else {
					if err := bot.StorageRepo.SaveLastSuccessfulUpdateID(update.UpdateID); err != nil {
						errorChan <- err
					}
				}
			}(update)
		}
	}
}
