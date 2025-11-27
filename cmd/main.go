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

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	aiProviderToken := os.Getenv("AI_TOKEN")

	storageRepo := telegrambot.NewInmemoryBotStorageRepository()
	mentorRepo := ai.NewAIMentorRepository(aiProviderToken)
	bot, err := telegrambot.NewBot(telegramBotToken, whitelist, storageRepo, mentorRepo)
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
	whitelistStrings := strings.Split(whitelistCleaned, ",")
	for _, s := range whitelistStrings {
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

	log.Printf("Starting bot with offset: %d", offset)

	for {
		updates, err := bot.GetUpdates(offset)
		if err != nil {
			return err
		}
		for _, update := range updates {
			log.Printf("Received update: %+v", update)
			if update.Message != nil {
				telegrambot.BaseMiddleware(bot, update.Message)
			}
			offset = update.UpdateID + 1
			if err := bot.StorageRepo.SaveLastSuccessfulUpdateID(update.UpdateID); err != nil {
				return err
			}
		}
	}
}
