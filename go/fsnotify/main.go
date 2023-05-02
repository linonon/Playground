package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	filePath := "track.json"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		// 获取文件的初始大小
		fi, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
		}
		oldSize := fi.Size()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					// 获取文件的当前大小
					fi, err := os.Stat(filePath)
					if err != nil {
						log.Fatal(err)
					}
					newSize := fi.Size()

					// 读取新写入的内容
					f, err := os.Open(filePath)
					if err != nil {
						log.Fatal(err)
					}
					defer f.Close()

					buf := make([]byte, newSize-oldSize)
					_, err = f.ReadAt(buf, oldSize)
					if err != nil {
						log.Fatal(err)
					}

					// 更新文件大小
					oldSize = newSize

					// 在这里解码 JSON 并根据 "type" 字段执行 MongoDB 存储操作
					fmt.Printf("new content: %s\n", buf)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("error: %v\n", err)
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
