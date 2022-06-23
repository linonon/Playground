package lo_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
)

func TestAttempt(t *testing.T) {
	{
		iter, err := lo.Attempt(42, func(i int) error {
			if i == 5 {
				return nil
			}

			return fmt.Errorf("failed")
		})

		fmt.Println(iter, err)
	}
	{
		iter, err := lo.Attempt(2, func(i int) error {
			if i == 2 {
				return nil
			}

			return fmt.Errorf("failed")
		})

		fmt.Println(iter, err)
	}
	{
		iter, err := lo.Attempt(0, func(i int) error {
			if i < 42 {
				return fmt.Errorf("failed")
			}

			return nil
		})

		fmt.Println(iter, err)
	}
}

func TestAttemptWithDelay(t *testing.T) {
	iter, duration, err := lo.AttemptWithDelay(5, 2*time.Second, func(i int, duration time.Duration) error {
		if i == 2 {
			return nil
		}

		return fmt.Errorf("failed")
	})

	fmt.Println(iter, duration, err)
}

func TestDebounce(t *testing.T) {
	f := func() {
		fmt.Println("Call once after 1500ms when debounce stopped invoking!")
	}

	reset, cancel := lo.NewDebounce(1500*time.Millisecond, f)
	for j := 0; j < 10; j++ {
		reset()
	}
	now := time.Now()
	fmt.Println(time.Now().Sub(now))

	time.Sleep(1 * time.Second)
	fmt.Println(time.Now().Sub(now))
	// cancel() // will not print f()

	time.Sleep(1 * time.Second)

	fmt.Println(time.Now().Sub(now))

	cancel()
}

func TestSynchronize(t *testing.T) {
	s := lo.Synchronize()

	for i := 0; i < 10; i++ {
		go s.Do(func() {
			println("will be called sequentially")
		})
	}

	time.Sleep(time.Second)
}

func TestAsync(t *testing.T) {
	{
		ch := lo.Async(func() error { time.Sleep(3 * time.Second); return nil })
		fmt.Println(<-ch)
	}
	{
		ch := lo.Async0(func() { time.Sleep(3 * time.Second) })
		fmt.Println(<-ch)
	}
	{
		ch := lo.Async1(func() int { time.Sleep(3 * time.Second); return 42 })
		fmt.Println(<-ch)
	}
	{
		ch := lo.Async2(func() (int, string) { time.Sleep(3 * time.Second); return 42, "hello" })
		fmt.Println(<-ch)
	}
}
