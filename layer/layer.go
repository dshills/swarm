package layer

import (
	"log"
	"sync"

	"github.com/dshills/swarm/comm"
	"github.com/dshills/swarm/def"
	"github.com/dshills/swarm/lua"
	"github.com/dshills/swarm/neuron"
	"github.com/dshills/swarm/signal"
	"github.com/dshills/swarm/task"
	"github.com/dshills/swarm/trans"
	"github.com/google/uuid"
)

type Layer interface {
	Describe() string
	Consider(signal.Signal, task.List) signal.Signal
}

func New(ld def.LayerDef, com *comm.Comms, luaFns []lua.Function) Layer {
	l := _layer{id: uuid.New().String(), ignoreContext: ld.IgnoreContext, luaFns: luaFns}
	for _, mod := range ld.NeuronModels {
		n := neuron.New(ld, com, mod)
		trans := trans.New(n, ld, luaFns)
		l.transmitters = append(l.transmitters, trans)
		l.neurons = append(l.neurons, n)
	}
	return &l
}

type _layer struct {
	neurons       []neuron.Neuron
	transmitters  []trans.Transmitter
	id            string
	ignoreContext bool
	luaFns        []lua.Function
}

func (l *_layer) Describe() string {
	return ""
}

func (l *_layer) Consider(s signal.Signal, tl task.List) signal.Signal {
	if tl.Len() == 0 {
		return s
	}
	l.distributeTask(s, tl)
	return l.updateTasks(s, tl)
}

func (l *_layer) distributeTask(s signal.Signal, tl task.List) {
	wg := sync.WaitGroup{}
	for _, trans := range l.transmitters {
		task := tl.Pop()
		if task == "" {
			break
		}
		newSig := s.NewChild(task, l.ignoreContext)
		wg.Add(1)
		go l.transmit(&wg, trans, s, newSig)
	}
	wg.Wait()
	if tl.Len() > 0 {
		l.distributeTask(s, tl)
	}
}

func (l *_layer) transmit(wg *sync.WaitGroup, t trans.Transmitter, org, s signal.Signal) {
	defer wg.Done()
	s = t.Transmit(s)
	if s.Err() != nil {
		log.Printf("[ERROR] %v %v", s.ModelUsed(), s.Err())
		org.SetError(s.Err())
	}
}

func (l *_layer) updateTasks(s signal.Signal, tl task.List) signal.Signal {
	if len(s.ParsedResult()) > 0 {
		tl.Push(s.ParsedResult()...)
		return s
	}
	tl.Push(s.CombinedResult()...)
	return s
}
