package core

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sync"
)

type State struct {
	sync.RWMutex
	processed map[string]bool
}

func NewState() *State {
	return &State{
		processed: make(map[string]bool),
	}
}

func (s *State) makeTaskId(data string, agentId string) string {
	h := sha1.New()
	io.WriteString(h, data)
	io.WriteString(h, agentId)
	return fmt.Sprintf("% x", h.Sum(nil))
}

func (s *State) Add(data string, agentId string) {
	s.Lock()
	defer s.Unlock()
	s.processed[s.makeTaskId(data, agentId)] = true
}

func (s *State) DidProcess(data string, agentId string) (found bool) {
	s.RLock()
	defer s.RUnlock()
	_, found = s.processed[s.makeTaskId(data, agentId)]
	return
}
