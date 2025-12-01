package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	whitelistEnvKey       = "WHITELIST_IDS"
	telegramBotEnvKey     = "TELEGRAM_BOT_TOKEN"
	aiProviderTokenEnvKey = "AI_TOKEN"
	adminIdEnvKey         = "ADMIN_ID"
)

const (
	telegramBaseUrl = "https://api.telegram.org/"
	genApiBaseUrl   = "https://api.gen-api.ru/"
)

const systemPrompt = `You are a native English speaker in a casual English-learning chat. Sound natural and friendly, not academic. Correct language like a real native would say it.
Rules:
- The field <ExtraQuestion>= is optional.
- If <ExtraQuestion>= exists → answer ONLY that question (must be about grammar, spelling, vocabulary, etc.).
- If <ExtraQuestion>= is empty or missing → do NOT answer anything, only correct the text from <Question>= for grammar, spelling, and natural wording.
- Ignore the meaning of <Question>=.
- Keep responses short: 30–100 words, max 2000 characters.
- Output = plain text. Lists allowed: - or 1. 2. (symbols don’t count as words).
- Provide natural fixes + brief explanation if useful.`

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	whitelist, err := whitelist()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load whitelist: %w", err)
	}

	adminId, err := adminId()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load admin ID: %w", err)
	}

	telegramBotToken, err := telegramBotToken()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load Telegram bot token: %w", err)
	}

	aiProviderToken, err := aiProviderToken()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load AI provider token: %w", err)
	}

	config := Config{
		WhitelistIDs:         whitelist,
		AdminID:              adminId,
		TelegramBotToken:     telegramBotToken,
		AIProviderToken:      aiProviderToken,
		TelegramAPIBaseURL:   telegramBaseUrl,
		AIProviderAPIBaseURL: genApiBaseUrl,
		MentorSystemPrompt:   systemPrompt,
	}

	return config, nil
}

func whitelist() (map[int64]struct{}, error) {
	v := os.Getenv(whitelistEnvKey)
	v = strings.TrimSpace(v)

	if len(v) == 0 {
		return nil, fmt.Errorf("whitelist IDs cannot be empty")
	}

	const defaultCapacity = 10
	whitelist := make(map[int64]struct{}, defaultCapacity)

	stringIds := strings.SplitSeq(v, ",")
	for s := range stringIds {
		n, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			whitelist[n] = struct{}{}
		} else {
			slog.Warn("Failed to parse whitelist ID, skipping", "id", s, "error", err)
		}
	}
	return whitelist, nil
}

func adminId() (int64, error) {
	v := os.Getenv(adminIdEnvKey)
	v = strings.TrimSpace(v)

	if len(v) == 0 {
		return 0, fmt.Errorf("admin ID cannot be empty")
	}

	adminId, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse admin ID: %w", err)
	}

	return adminId, nil
}

func telegramBotToken() (string, error) {
	v := os.Getenv(telegramBotEnvKey)
	v = strings.TrimSpace(v)

	if len(v) == 0 {
		return "", fmt.Errorf("telegram bot token cannot be empty")
	}

	return v, nil
}

func aiProviderToken() (string, error) {
	v := os.Getenv(aiProviderTokenEnvKey)
	v = strings.TrimSpace(v)

	if len(v) == 0 {
		return "", fmt.Errorf("AI provider token cannot be empty")
	}

	return v, nil
}
