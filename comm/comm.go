package comm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	OpenAI   = "openai"
	Ollama   = "ollama"
	ollamaEP = "api/chat"
	openaiEP = "/chat/completions"
)

const (
	RoleSystem    = "system"
	RoleAssistant = "assistant"
	RoleUser      = "user"
)

type Comms struct {
	client http.Client
	models []model
}

func (c *Comms) Configure(confPath string) error {
	file, err := os.Open(confPath)
	if err != nil {
		return err
	}
	mods := []model{}
	if err := yaml.NewDecoder(file).Decode(&mods); err != nil {
		return err
	}
	c.models = mods
	return nil
}

func (c *Comms) Complete(model string, prompt string) (*Response, error) {
	msg := Message{Role: RoleUser, Content: prompt}
	return c.Converse(model, []Message{msg})
}

func (c *Comms) Converse(model string, msgs []Message) (*Response, error) {
	start := time.Now()
	mod, err := c.modelInfo(model)
	if err != nil {
		return nil, err
	}
	req := request{
		Model:    mod.Model,
		Messages: msgs,
		Stream:   false,
	}
	ep := ""
	if strings.EqualFold(mod.API, Ollama) {
		ep, err = url.JoinPath(mod.BaseURL, ollamaEP)
	} else if strings.EqualFold(mod.API, OpenAI) {
		ep, err = url.JoinPath(mod.BaseURL, openaiEP)
	}
	if err != nil {
		return nil, err
	}
	fmt.Println(ep)

	byts, _ := json.Marshal(&req)
	byts, err = c.send(ep, bytes.NewReader(byts), mod.APIKey)
	if err != nil {
		return nil, err
	}

	if strings.EqualFold(mod.API, OpenAI) {
		return c.openai(byts, start)
	}
	if strings.EqualFold(mod.API, Ollama) {
		return c.ollama(byts, start)
	}
	return nil, fmt.Errorf("unknown API %v", mod.API)
}

func (c *Comms) ollama(byts []byte, start time.Time) (*Response, error) {
	or := ollamaResponse{}
	if err := json.Unmarshal(byts, &or); err != nil {
		return nil, err
	}
	resp := or.AsResponse(time.Since(start))
	return &resp, nil
}

func (c *Comms) openai(byts []byte, start time.Time) (*Response, error) {
	or := openaiResponse{}
	if err := json.Unmarshal(byts, &or); err != nil {
		return nil, err
	}
	resp := or.AsResponse(time.Since(start))
	return &resp, nil
}

func (c *Comms) modelInfo(name string) (*model, error) {
	for _, m := range c.models {
		if strings.EqualFold(m.Model, name) {
			return &m, nil
		}
		for _, a := range m.Aliases {
			if strings.EqualFold(a, name) {
				return &m, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

func (c *Comms) send(endpoint string, reader io.Reader, apiKey string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("send: %v %v", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}
