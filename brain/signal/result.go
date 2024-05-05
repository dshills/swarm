package signal

import (
	"fmt"
	"strings"
	"time"
)

type Result struct {
	Completed bool
	Error     error
	Results   []string
	Usage     Usage
	Duration  Duration
}

type Duration struct {
	Actual time.Duration
	AI     time.Duration
}

type Usage struct {
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
}

func newResult(results []string, usage Usage, dur Duration, completed bool, err error) Result {
	return Result{
		Results:   results,
		Completed: completed,
		Error:     err,
		Usage:     usage,
		Duration:  dur,
	}
}

func (r *Result) String() string {
	builder := strings.Builder{}
	if r.Error != nil {
		builder.WriteString(fmt.Sprintf("ERROR: %v\n", r.Error))
	} else {
		builder.WriteString(fmt.Sprintf("Completed: %v\n", r.Completed))
	}
	builder.WriteString(fmt.Sprintf("AI Elapsed: %v\n", r.Duration.AI))
	builder.WriteString(fmt.Sprintf("Total Elapsed: %v\n", r.Duration.Actual))
	builder.WriteString(fmt.Sprintf("Prompt Tokens: %v\n", r.Usage.PromptTokens))
	builder.WriteString(fmt.Sprintf("Completion Tokens: %v\n", r.Usage.CompletionTokens))
	builder.WriteString(fmt.Sprintf("Total Tokens: %v\n", r.Usage.TotalTokens))
	return builder.String()
}
