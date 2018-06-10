package core

import (
	"fmt"
	"log"
	"sync"
)

type Orchestrator struct {
	sync.RWMutex

	dataBus *DataBus
	agents  map[string]Agent
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		dataBus: NewDataBus(),
		agents:  make(map[string]Agent),
	}
}

func (o *Orchestrator) Register(agent Agent) error {
	o.Lock()
	defer o.Unlock()

	if _, found := o.agents[agent.ID()]; found {
		return fmt.Errorf("an agent with id '%s' was already registered", agent.ID())
	}

	if err := agent.Register(o); err != nil {
		return fmt.Errorf("could not register agent '%s': %v", agent.ID(), err)
	}

	o.agents[agent.ID()] = agent

	return nil
}

func (o *Orchestrator) Publish(eventName string, args ...interface{}) {
	log.Printf("publish: %s(%v)", eventName, args)
	o.dataBus.Publish(eventName, args...)
}

func (o *Orchestrator) Subscribe(eventName string, fn interface{}) {
	o.dataBus.SubscribeAsync(eventName, fn, false)
}

func (o *Orchestrator) Wait() {
	o.dataBus.WaitAsync()
}
