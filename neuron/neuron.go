package neuron

import (
	"log"
	"strings"
	"time"

	"github.com/dshills/swarm/comm"
	"github.com/dshills/swarm/def"
	"github.com/dshills/swarm/signal"
	"github.com/google/uuid"
)

type Neuron interface {
	Describe() string
	Work(signal.Signal) signal.Signal
}

func New(ld def.LayerDef, com *comm.Comms, modName string) Neuron {
	return &_neuron{id: uuid.New().String(), com: com, layerDef: ld, modelName: modName}
}

const taskInsert = "%%TASK%%"

type _neuron struct {
	id        string
	com       *comm.Comms
	layerDef  def.LayerDef
	modelName string
	conv      []comm.Message
}

func (n _neuron) Describe() string {
	return n.id
}

func (n _neuron) Work(s signal.Signal) signal.Signal {
	log.Printf("Starting Neuron %v with %v\n", n.id, n.modelName)
	start := time.Now()
	n.makeConv(s)

	resp, err := n.com.Converse(n.modelName, n.conv)
	if err != nil {
		s.SetError(err)
		return s
	}

	n.conv = append(n.conv, resp.Message)
	s.SetResult(resp.Message.Content, n.modelName)

	s.SetElapsed(time.Since(start), resp.Elapsed)
	s.SetTokens(resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	log.Printf("Completed Neuron %v\n", n.id)
	return s
}

func (n *_neuron) makeConv(s signal.Signal) {
	prompt := n.layerDef.Prompt
	if strings.Contains(n.layerDef.Prompt, taskInsert) {
		prompt = strings.ReplaceAll(n.layerDef.Prompt, taskInsert, s.Task())
	} else {
		prompt += " " + s.Task()
	}
	msgPrompt := comm.Message{Role: comm.RoleUser, Content: prompt}
	sys := comm.Message{Role: comm.RoleSystem, Content: n.layerDef.Persona}

	s.AddToConversation(sys, msgPrompt)
	// Add our messages to NewConv
	n.conv = append(n.conv, s.PrevConversation()...)
	n.conv = append(n.conv, sys)
	n.conv = append(n.conv, msgPrompt)
}
