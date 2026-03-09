package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var (
	x   int
	m   sync.Mutex
	wg  sync.WaitGroup
	sum int
)

func Add() {
	defer wg.Done()
	data := rand.Intn(10000)
	sum += data
	for i := 0; i < data; i++ {
		m.Lock()
		x += 1
		m.Unlock()
	}
}

func main() {
	num := rand.Intn(31) + 1
	wg.Add(num)
	for i := 0; i < num; i++ {
		go Add()
	}
	wg.Wait()
	fmt.Println("——————计算完成！——————")
	fmt.Printf("理论上的和为: %d ;实际上的和为: %d\n", sum, x)
}
