package comm

import "time"

type model struct {
	Host       string   `yaml:"Host"`
	Model      string   `yaml:"Model"`
	API        string   `yaml:"API"`
	BaseURL    string   `yaml:"BaseURL"`
	APIKey     string   `yaml:"APIKey"`
	Promptcost float64  `yaml:"PromptCost"`
	Compcost   float64  `yaml:"CompCost"`
	Aliases    []string `yaml:"Aliases"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Message Message
	Usage   Usage
	Elapsed time.Duration
}

type Usage struct {
	CompletionTokens int64 `json:"completion_tokens"`
	PromptTokens     int64 `json:"prompt_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ollamaResponse struct {
	Model           string  `json:"model"`
	Message         Message `json:"message"`
	Done            bool    `json:"done"`
	PromptEvalCount int64   `json:"prompt_eval_count"`
	EvalCount       int64   `json:"eval_count"`
}

func (or *ollamaResponse) AsResponse(elapsed time.Duration) Response {
	us := Usage{
		PromptTokens: or.PromptEvalCount,
		TotalTokens:  or.EvalCount,
	}
	return Response{
		Message: or.Message,
		Usage:   us,
		Elapsed: elapsed,
	}
}

type openaiResponse struct {
	Choices []openaiChoice `json:"choices,omitempty"`
	Usage   Usage          `json:"usage,omitempty"`
}

func (oa *openaiResponse) AsResponse(elapsed time.Duration) Response {
	return Response{
		Message: oa.Choices[0].Message,
		Elapsed: elapsed,
		Usage:   oa.Usage,
	}
}

type openaiChoice struct {
	FinishReason string  `json:"finish_reason,omitempty"`
	Message      Message `json:"message,omitempty"`
}
