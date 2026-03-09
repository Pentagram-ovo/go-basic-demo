package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Work(ctx context.Context, id int) {
	defer wg.Done()
	fmt.Printf("Worker %d 开始工作！！！\n", id)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("检测到程序退出，Worker %d 爽爽停止工作！\n", id)
			return
		default:
			fmt.Printf("Worker %d 工作ing！！！\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	//手动退出
	//ctx, cancel := context.WithCancel(context.Background())
	//超时退出(超过三秒自动退出)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	workernum := rand.Intn(10) + 1
	wg.Add(workernum)
	for i := 0; i < workernum; i++ {
		go Work(ctx, i)
	}
	//手动需要
	//time.Sleep(time.Second)
	//fmt.Println("准备下班拉闸！")
	//超时设置就不用手动调用~
	//cancel()

	wg.Wait()
	fmt.Println("————所有工人已下班，工厂关门！————")
}
