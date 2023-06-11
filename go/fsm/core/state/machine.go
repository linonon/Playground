package state

import "time"

type State interface {
	// Name returns the name of the state.
	// This is used to identify the state in the StateMachine.
	// The name must be unique.
	Name() string
	// OnEntry is called when the state is entered.
	OnEntry()
	// OnExit is called when the state is exited.
	OnExit()
	// TimeoutDuration returns the duration of the state's timeout.
	TimeoutDuration() time.Duration
	// TimeoutFunc is called when the state's timeout expires.
	TimeoutFunc()
	// Next returns the name of the next state.
	// This is used to transition to the next state.
	// The name must be unique.
	Next() string
}

type StateMachine struct {
	States  map[string]State
	Current State
	Timer   *time.Timer
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		States: make(map[string]State),
	}
}

// AddStates adds a state to the StateMachine.
func (sm *StateMachine) AddStates(states ...State) *StateMachine {
	for _, s := range states {
		sm.States[s.Name()] = s
	}
	return sm
}

// SetCurrent sets the current state of the StateMachine.
func (sm *StateMachine) SetStartState(s string) *StateMachine {
	return sm.SetCurrent(s)
}

// SetCurrent sets the current state of the StateMachine.
func (sm *StateMachine) SetCurrent(s string) *StateMachine {
	if sm.Current != nil {
		sm.Current.OnExit()
	}
	sm.Current = sm.States[s]
	sm.Current.OnEntry()
	sm.Timer = time.AfterFunc(sm.Current.TimeoutDuration(), sm.Current.TimeoutFunc)

	return sm
}

// Transition transitions to the specified state.
func (sm *StateMachine) Transition(s string) {
	if sm.Timer != nil {
		sm.Timer.Stop()
	}
	sm.SetCurrent(s)
}

// DefaultTransition transitions to the next state.
func (sm *StateMachine) DefaultTransition() {
	if sm.Timer != nil {
		sm.Timer.Stop()
	}
	sm.SetCurrent(sm.Current.Next())
}
