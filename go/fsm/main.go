package main

import (
	"fmt"
	"math/rand"
	"myfsm/core/state"
	"time"
)

const (
	StateA = "a"
	StateB = "b"
	StateC = "c"
	StateD = "d"
)

func main() {

	sm := state.NewStateMachine().
		AddStates(
			state.NewDefaultState(StateA).
				SetNextFunc(func() string { return "b" }).
				SetEntryFunc(func() { fmt.Println("Entering state a") }).
				SetExitFunc(func() { fmt.Println("Exiting state a") }).
				SetTimeoutFunc(func() { fmt.Println("Timeout in state a") }).
				SetTimeoutDuration(2*time.Second),

			state.NewDefaultState(StateB).
				SetNextFunc(func() string {
					if rand.Intn(2) == 0 {
						return "c"
					} else {
						return "d"
					}
				}).
				SetEntryFunc(func() { fmt.Println("Entering state b") }).
				SetExitFunc(func() { fmt.Println("Exiting state b") }).
				SetTimeoutFunc(func() { fmt.Println("Timeout in state b") }).
				SetTimeoutDuration(2*time.Second),

			state.NewDefaultState(StateC).
				SetNextFunc(func() string { return "a" }).
				SetEntryFunc(func() { fmt.Println("Entering state c") }).
				SetExitFunc(func() { fmt.Println("Exiting state c") }).
				SetTimeoutFunc(func() { fmt.Println("Timeout in state c") }).
				SetTimeoutDuration(2*time.Second),

			state.NewDefaultState(StateD).
				SetNextFunc(func() string { return "a" }).
				SetEntryFunc(func() { fmt.Println("Entering state d") }).
				SetExitFunc(func() { fmt.Println("Exiting state d") }).
				SetTimeoutFunc(func() { fmt.Println("Timeout in state d") }).
				SetTimeoutDuration(2*time.Second),
		).
		SetStartState(StateA)

	for {
		time.Sleep(3 * time.Second)
		sm.DefaultTransition()
	}
}
