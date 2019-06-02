package goteam

import (
	"errors"
	"sync"
)

type Worker interface {
	Perform(t Task) error
}

type Task struct {
	Data string
}

type Team struct {
	ch chan Worker
	wg sync.WaitGroup
}

func New(capacity int) (*Team, error) {
	if capacity < 1 {
		return nil, errors.New("nonpositive team capacity")
	}
	team := Team{
		ch: make(chan Worker, capacity),
	}
	return &team, nil
}

func (team *Team) Add(w Worker) error {
	select {
	case team.ch <- w:
		return nil
	default:
		return errors.New("team capacity reached")
	}
}

func (team *Team) Accept(task Task) {
	team.wg.Add(1)
	go func(Task) {
		defer team.wg.Done()
		worker := <-team.ch
		worker.Perform(task) // error ignored for now
		team.ch <- worker
	}(task)
}

func (team *Team) Shutdown() {
	team.wg.Wait()
}
