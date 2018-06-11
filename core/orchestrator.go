package core

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

type Task func()

type Orchestrator struct {
	sync.RWMutex

	workers int
	tasks   chan Task
	wg      sync.WaitGroup

	dataBus *DataBus
	state   *State
	agents  map[string]Agent
}

func NewOrchestrator(workers int) *Orchestrator {
	if workers <= 0 {
		workers = runtime.NumCPU() * 2
	}

	return &Orchestrator{
		workers: workers,
		tasks:   make(chan Task),
		wg:      sync.WaitGroup{},
		dataBus: NewDataBus(),
		state:   NewState(),
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

func (o *Orchestrator) worker(id int) {
	// log.Printf("started worker #%d", id)

	for task := range o.tasks {
		if task == nil {
			log.Printf("stopping worker %d", id)
			return
		}

		// log.Printf("running task %v", task)
		task()

		o.wg.Done()
	}
}

func (o *Orchestrator) Start() {
	log.Printf("starting %d workers...", o.workers)

	for i := 0; i < o.workers; i++ {
		go o.worker(i)
	}
}

func (o *Orchestrator) RunTask(t Task) {
	o.wg.Add(1)
	o.tasks <- t
}

func (o *Orchestrator) Publish(eventName string, args ...interface{}) {
	key := fmt.Sprintf("%s(%v)", eventName, args)

	if o.state.DidProcess(key, "main:state") {
		return
	}

	o.state.Add(key, "main:state")

	log.Printf("publish: \033[1m%s\033[0m", key)

	o.dataBus.Publish(eventName, args...)
}

func (o *Orchestrator) Subscribe(eventName string, fn interface{}) {
	o.dataBus.SubscribeAsync(eventName, fn, false)
}

func (o *Orchestrator) Wait() {
	o.dataBus.WaitAsync()
	o.wg.Wait()
}
