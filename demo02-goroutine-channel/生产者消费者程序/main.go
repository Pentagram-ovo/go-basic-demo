package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	//分开定义防止deadlock死锁
	var prowg sync.WaitGroup
	var conwg sync.WaitGroup
	//定义通道并设置缓冲
	ch := make(chan int, 3)
	//随机生成数量
	producer_num := rand.Intn(5) + 1
	consumer_num := rand.Intn(10) + 1

	fmt.Printf("————共有 %d 个生产者，共有 %d 个消费者————\n", producer_num, consumer_num)
	//消费者
	for i := 0; i < consumer_num; i++ {
		conwg.Add(1)
		go func(consumerId int) {
			defer conwg.Done()
			//循环读取通道内数据，直到通道关闭
			for data := range ch {
				fmt.Printf("第 %d 个消费者消费了 %d 的数据\n", consumerId+1, data)
			}
			//模拟消费所需时间
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}(i)
	}

	//生产者
	for i := 0; i < producer_num; i++ {
		prowg.Add(1)
		go func(producerId int) {
			defer prowg.Done()
			sum := 0
			//假设每个生产者生产7条数据
			for j := 0; j < 7; j++ {
				data := rand.Intn(1000)
				ch <- data
				sum += data
			}
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			fmt.Printf("第 %d 个生产者共生产了 %d 的数据\n", producerId+1, sum)
		}(i)
	}
	//等待生产者全部生产完后再关闭通道
	go func() {
		prowg.Wait()
		close(ch)
	}()
	conwg.Wait()

	fmt.Println("——————————进程全部完成~——————————")
}
