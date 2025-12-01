package ai

type GenApiRequest struct {
	Messages         []InputMessage     `json:"messages"`
	CallbackUrl      string             `json:"callback_url,omitempty"`
	IsSync           bool               `json:"is_sync,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	N                int                `json:"n,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	ResponseFormat   string             `json:"response_format,omitempty"`
	Tools            []string           `json:"tools,omitempty"`
	ToolChoice       string             `json:"tool_choice,omitempty"`
	Seed             int64              `json:"seed,omitempty"`
}

type InputMessage struct {
	Role    string        `json:"role"`
	Content []ContentItem `json:"content"`
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type GenApiResponse struct {
	RequestId int64      `json:"request_id"`
	Model     string     `json:"model"`
	Cost      float64    `json:"cost"`
	Responses []Response `json:"response"`
}

type Response struct {
	Index        int           `json:"index"`
	Message      MessageOutput `json:"message"`
	FinishReason string        `json:"finish_reason,omitempty"`
}

type MessageOutput struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Refusal bool   `json:"refusal,omitempty"`
}
