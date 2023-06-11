package state

import "time"

type DefaultState struct {
	name            string
	nextFunc        func() string
	entryFunc       func()
	exitFunc        func()
	timeoutFunc     func()
	timeoutDuration time.Duration

	err error
}

// NewDefaultState returns a new DefaultState.
func NewDefaultState(name string) *DefaultState {
	return &DefaultState{
		name: name,
	}
}

//////////////////////////////////////////////
// Implement State interface
//////////////////////////////////////////////

func (s *DefaultState) Name() string {
	return s.name
}

func (s *DefaultState) OnEntry() {
	if s.entryFunc != nil {
		s.entryFunc()
	}
}

func (s *DefaultState) OnExit() {
	if s.exitFunc != nil {
		s.exitFunc()
	}
}

func (s *DefaultState) TimeoutFunc() {
	if s.timeoutFunc != nil {
		s.timeoutFunc()
	}
}

func (s *DefaultState) TimeoutDuration() time.Duration {
	return s.timeoutDuration
}

func (s *DefaultState) Next() string {
	return s.nextFunc()
}

//////////////////////////////////////////////
// Setters
//////////////////////////////////////////////

func (s *DefaultState) SetNextFunc(nextFunc func() string) *DefaultState {
	s.nextFunc = nextFunc
	return s
}

func (s *DefaultState) SetEntryFunc(entryFunc func()) *DefaultState {
	s.entryFunc = entryFunc
	return s
}

func (s *DefaultState) SetExitFunc(exitFunc func()) *DefaultState {
	s.exitFunc = exitFunc
	return s
}

func (s *DefaultState) SetTimeoutFunc(timeoutFunc func()) *DefaultState {
	s.timeoutFunc = timeoutFunc
	return s
}

func (s *DefaultState) SetTimeoutDuration(timeoutDuration time.Duration) *DefaultState {
	s.timeoutDuration = timeoutDuration
	return s
}
