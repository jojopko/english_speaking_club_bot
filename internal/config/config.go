package config

type Config struct {
	WhitelistIDs         map[int64]struct{}
	AdminID              int64
	TelegramBotToken     string
	AIProviderToken      string
	TelegramAPIBaseURL   string
	AIProviderAPIBaseURL string
	MentorSystemPrompt   string
	TelegramBotNickname  string
}
