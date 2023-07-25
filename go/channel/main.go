package main

import (
	"fmt"
	"sync"
	"time"
)

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type testChan[T Int] chan T

func NewTestChan[T Int](size int) testChan[T] {
	return make(chan T, size)
}

func (ch testChan[T]) Send(i T) {
	ch <- i
}

func (ch testChan[T]) Recive() (T, error) {
	if ch == nil {
		return T(-1), fmt.Errorf("channel is nil")
	}

	r := <-ch
	if r == -1 {
		ch.Close()
		return -1, fmt.Errorf("channel is closed")
	}

	return r, nil
}

func (ch testChan[T]) Close() {
	close(ch)
}

func main() {
	wg := sync.WaitGroup{}

	var ch = NewTestChan[int](1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			ch.Send(1)
			time.Sleep(1 * time.Second)
		}

		ch.Send(-1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			i, err := ch.Recive()
			if err != nil {
				break
			}

			fmt.Println(i)
		}
	}()

	wg.Wait()
	fmt.Println("Done")
}
