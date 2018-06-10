package core

type Agent interface {
	ID() string
	Register(o *Orchestrator) error
}
