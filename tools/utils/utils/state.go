package utils

import (
	"sync"
	"time"
)

const (
	StateClosed    = -1
	StateDeadline  = -2
	StateWaitClose = -3
	StateTimeout   = 1
)

func NewState() *State {
	return &State{
		base: &stateBase{
			done: make(chan struct{}),
		},
	}
}

type stateBase struct {
	wg     sync.WaitGroup
	done   chan struct{}
	mux    sync.Mutex
	stat   int
	timer  *time.Timer
	errors []error
}

type State struct {
	base  *stateBase
	stat  int
	timer *time.Timer
}

func (s *State) New(d time.Duration) *State {
	s.base.wg.Add(1)
	state := &State{base: s.base}
	if d > 0 {
		state.SetTimeout(d)
	}
	return state
}

func (s *State) Delete() {
	s.base.wg.Done()
}

func (s *State) SetDeadline(d time.Duration, closed bool) {
	s.base.mux.Lock()
	defer s.base.mux.Unlock()
	if s.base.stat == 0 {
		if closed {
			s.base.stat = StateWaitClose
		}
		if s.base.timer == nil {
			s.base.timer = time.AfterFunc(d, func() {
				s.base.mux.Lock()
				defer s.base.mux.Unlock()
				if s.base.stat == 0 || s.base.stat == StateWaitClose {
					s.base.stat = StateDeadline
					close(s.base.done)
				}
			})
		} else {
			s.base.timer.Reset(d)
		}
	}
}

func (s *State) SetTimeout(d time.Duration) {
	s.base.mux.Lock()
	defer s.base.mux.Unlock()
	if s.base.stat == 0 {
		s.stat = 0
		if s.timer == nil {
			c := make(chan time.Time, 1)
			s.timer = time.AfterFunc(d, func() {
				s.base.mux.Lock()
				defer s.base.mux.Unlock()
				s.stat = StateTimeout
				select {
				case c <- time.Now():
				default:
				}
			})
			s.timer.C = c
		} else {
			s.timer.Reset(d)
		}
	}
}

func (s *State) Done() <-chan struct{} {
	return s.base.done
}

func (s *State) Timeout() <-chan time.Time {
	s.base.mux.Lock()
	defer s.base.mux.Unlock()
	if s.timer != nil {
		return s.timer.C
	}
	return nil
}

func (s *State) Get() int {
	if s.base.stat != 0 {
		return s.base.stat
	} else if s.stat != 0 {
		return s.stat
	} else if s.timer == nil {
		return 0
	}
	select {
	case <-s.base.done:
		return s.base.stat
	case <-s.timer.C:
		return s.stat
	}
}

func (s *State) SetError(e error) {
	s.base.mux.Lock()
	defer s.base.mux.Unlock()
	s.base.errors = append(s.base.errors, e)
}

func (s *State) Error() error {
	s.base.mux.Lock()
	defer s.base.mux.Unlock()
	if len(s.base.errors) > 0 {
		return s.base.errors[0]
	}
	return nil
}

func (s *State) Errors() []error {
	s.base.mux.Lock()
	defer s.base.mux.Unlock()
	return s.base.errors
}

func (s *State) SetClosed() {
	s.base.mux.Lock()
	defer s.base.mux.Unlock()
	if s.base.stat == 0 || s.base.stat == StateWaitClose {
		s.base.stat = StateClosed
		close(s.base.done)
	}
}

func (s *State) WaitClosed() {
	s.base.wg.Wait()
	s.SetClosed()
}
