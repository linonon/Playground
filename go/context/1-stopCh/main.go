package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	stopCh := make(chan bool)

	go func() {
		defer wg.Done()
		watchDog(stopCh, "監控狗")
	}()

	time.Sleep(5 * time.Second)
	stopCh <- true

	wg.Wait()
}

func watchDog(stopCh chan bool, name string) {
	for {
		select {
		case <-stopCh:
			fmt.Println(name, "有內鬼, 終止交易")
			return
		default:
			fmt.Println(name, "正在監控...")
		}

		time.Sleep(1 * time.Second)
	}
}
