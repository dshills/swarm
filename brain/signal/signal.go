package signal

import (
	"time"

	"github.com/dshills/swarm/brain/comm"
	"github.com/google/uuid"
)

type Signal interface {
	NewChild(task string, ignoreContext bool) Signal
	ID() string

	Err() error
	SetError(err error)

	Final() Result
	ModelUsed() string

	SetResult(result, modelUsed string)
	Result() string
	CombinedResult() []string

	ParsedResult() []string
	SetParsedResult(...string)

	SetComplete()
	IsComplete() bool

	AddToConversation(...comm.Message)
	SetPrevConversation(msgs ...comm.Message)
	PrevConversation() []comm.Message

	SetTask(string)
	Task() string

	SetElapsed(actual, ai time.Duration)
	SetTokens(req, cmp, total int64)
}

type _signal struct {
	id           string
	signals      []*_signal
	complete     bool
	err          error
	dur          Duration
	usage        Usage
	result       string
	parsedOutput []string
	task         string
	modelUsed    string
	prevConv     []comm.Message
	newConv      []comm.Message
}

func New(task string) Signal {
	return &_signal{task: task, id: uuid.New().String()}
}

func (s *_signal) NewChild(task string, ignoreContext bool) Signal {
	ns := &_signal{id: uuid.New().String(), task: task, prevConv: s.prevConv}
	if ignoreContext {
		ns = &_signal{id: uuid.New().String(), task: task}
	}
	s.signals = append(s.signals, ns)
	return ns
}

func (s *_signal) ID() string {
	return s.id
}

func (s *_signal) Err() error {
	return s.err
}

func (s *_signal) SetError(err error) {
	s.err = err
}

func (s *_signal) Final() Result {
	usage := s.usage
	dur := s.dur
	for _, sig := range s.signals {
		usage.PromptTokens += sig.usage.PromptTokens
		usage.CompletionTokens += sig.usage.CompletionTokens
		usage.TotalTokens += sig.usage.TotalTokens
		dur.AI += sig.dur.AI
		dur.Actual += sig.dur.Actual
	}
	return newResult(s.CombinedResult(), usage, dur, s.complete, s.err)
}

func (s *_signal) ModelUsed() string {
	return s.modelUsed
}

func (s *_signal) Result() string {
	return s.result
}

func (s *_signal) CombinedResult() []string {
	res := []string{s.result}
	for _, sig := range s.signals {
		res = append(res, sig.result)
	}
	return res
}

func (s *_signal) SetResult(result, modelUsed string) {
	s.modelUsed = modelUsed
	s.result = result
}

func (s *_signal) SetParsedResult(strs ...string) {
	s.parsedOutput = append(s.parsedOutput, strs...)
}

func (s *_signal) ParsedResult() []string {
	pres := s.parsedOutput
	for _, res := range s.signals {
		pres = append(pres, res.parsedOutput...)
	}
	return pres
}

func (s *_signal) SetComplete() {
	s.complete = true
}

func (s *_signal) IsComplete() bool {
	for _, sig := range s.signals {
		if sig.IsComplete() {
			return true
		}
	}
	return s.complete
}

func (s *_signal) AddToConversation(msgs ...comm.Message) {
	s.newConv = append(s.newConv, msgs...)
}

func (s *_signal) SetPrevConversation(msgs ...comm.Message) {
	s.prevConv = append(s.prevConv, msgs...)
}

func (s *_signal) PrevConversation() []comm.Message {
	return s.prevConv
}

func (s *_signal) Task() string {
	return s.task
}

func (s *_signal) SetTask(tsk string) {
	s.task = tsk
}

func (s *_signal) SetElapsed(actual, ai time.Duration) {
	s.dur.AI = ai
	s.dur.Actual = actual
}

func (s *_signal) SetTokens(req, cmp, total int64) {
	s.usage.PromptTokens = req
	s.usage.CompletionTokens = cmp
	s.usage.TotalTokens = total
}
