package trans

import (
	"fmt"

	"github.com/dshills/swarm/def"
	"github.com/dshills/swarm/lua"
	"github.com/dshills/swarm/neuron"
	"github.com/dshills/swarm/signal"
)

type Transmitter interface {
	Transmit(signal.Signal) signal.Signal
}

func New(to neuron.Neuron, ldef def.LayerDef, luaFns []lua.Function) Transmitter {
	return &_transmitter{to: to, ldef: ldef, luaFns: luaFns}
}

type _transmitter struct {
	to     neuron.Neuron
	ldef   def.LayerDef
	luaFns []lua.Function
}

func (t *_transmitter) Transmit(s signal.Signal) signal.Signal {
	// Update the task
	tsk, err := t.hookUpdateTask(s.Task())
	if err != nil {
		s.SetError(err)
		return s
	}
	s.SetTask(tsk)

	// Do the work
	s = t.to.Work(s)

	// Update the result
	res, err := t.hookUpdateResult(s.Result())
	if err != nil {
		s.SetError(err)
		return s
	}
	s.SetResult(res, s.ModelUsed())

	// Parse the result
	parsed, err := t.hookParseResults(s.Result())
	if err != nil {
		s.SetError(err)
		return s
	}
	s.SetParsedResult(parsed...)

	return s
}

func (t *_transmitter) fnLookup(name string) (lua.Function, error) {
	for _, fn := range t.luaFns {
		if fn.Name() == name {
			return fn, nil
		}
	}
	return nil, fmt.Errorf("lua function %v not found", name)
}

func (t *_transmitter) hookUpdateTask(task string) (string, error) {
	ntask := task
	for _, fname := range t.ldef.ChangeTaskFns {
		fn, err := t.fnLookup(fname)
		if err != nil {
			return task, err
		}
		ntask, err = fn.OneToOne(ntask)
		if err != nil {
			return task, err
		}
	}
	return ntask, nil
}

func (t *_transmitter) hookUpdateResult(result string) (string, error) {
	nresult := result
	for _, fname := range t.ldef.ChangeResultFns {
		fn, err := t.fnLookup(fname)
		if err != nil {
			return result, err
		}
		nresult, err = fn.OneToOne(nresult)
		if err != nil {
			return result, err
		}
	}
	return nresult, nil
}

func (t *_transmitter) hookParseResults(result string) ([]string, error) {
	parsed := []string{}
	for _, fname := range t.ldef.ResultToTasksFns {
		fn, err := t.fnLookup(fname)
		if err != nil {
			return nil, err
		}
		p, err := fn.OneToMany(result)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, p...)
	}
	return parsed, nil
}
