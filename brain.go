package brain

import (
	"fmt"
	"sync"

	"github.com/dshills/swarm/comm"
	"github.com/dshills/swarm/def"
	"github.com/dshills/swarm/layer"
	"github.com/dshills/swarm/lua"
	"github.com/dshills/swarm/signal"
	"github.com/dshills/swarm/task"
)

type Brain interface {
	Think(string) <-chan signal.Result
}

type _brain struct {
	thoughts []thought
	layers   []layer.Layer
	luaFuncs []lua.Function
}

type thought struct {
	wg       *sync.WaitGroup
	signal   signal.Signal
	resCh    chan<- signal.Result
	taskList task.List
}

func newBrain(bd *def.Definition, com *comm.Comms, luaFns []lua.Function) (Brain, error) {
	if len(bd.LayerDefs) == 0 {
		return nil, fmt.Errorf("layers must be > 0")
	}
	b := _brain{luaFuncs: luaFns}
	for _, ldef := range bd.LayerDefs {
		b.layers = append(b.layers, layer.New(ldef, com, luaFns))
	}
	return &b, nil
}

func (b *_brain) Think(taskVal string) <-chan signal.Result {
	s := signal.New(taskVal)
	ch := make(chan signal.Result)
	wg := sync.WaitGroup{}
	wg.Add(1)
	thought := thought{wg: &wg, signal: s, resCh: ch, taskList: task.New()}
	b.thoughts = append(b.thoughts, thought)
	go b.signal(thought, taskVal)
	return ch
}

func (b *_brain) Wait() {
	for _, t := range b.thoughts {
		t.wg.Wait()
	}
}

func (b *_brain) signal(t thought, task string) {
	defer t.wg.Done()
	t.taskList.Push(task)
	for _, l := range b.layers {
		s := l.Consider(t.signal, t.taskList)
		if s.IsComplete() || s.Err() != nil {
			t.resCh <- t.signal.Final()
		}
	}
	t.resCh <- t.signal.Final()
}
