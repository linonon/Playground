package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)
	ctx, stop := context.WithCancel(context.Background())
	valCtx := context.WithValue(ctx, "userID", 1006)

	go func() {
		defer wg.Done()
		getUser(valCtx)
	}()

	go func() {
		defer wg.Done()
		watchDog(ctx, "監控狗1")
	}()
	go func() {
		defer wg.Done()
		watchDog(ctx, "監控狗2")
	}()
	go func() {
		defer wg.Done()
		watchDog(ctx, "監控狗3")
	}()

	time.Sleep(5 * time.Second)
	stop()

	wg.Wait()
}

func watchDog(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "有內鬼, 終止交易")
			return
		case <-time.After(3 * time.Second):
			fmt.Println("超時!")
			return
		default:
			fmt.Println(name, "正在監控...")
		}

		time.Sleep(1 * time.Second)
	}
}

func getUser(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("獲取用戶, 攜程退出")
			return
		default:
			userID := ctx.Value("userID")
			fmt.Println("獲取用戶", "用戶 ID 為:", userID)
			time.Sleep(1 * time.Second)
		}
	}
}
