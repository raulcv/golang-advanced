package main

import (
	"fmt"
	"sync"
)

/*
Problem withdraw and deposit
go run --race build main.go -> to create file with information about concurrent status ok or not (Race condition)
*/
var (
	balance int = 100
)

func Deposit(amount int, wg *sync.WaitGroup, lock *sync.RWMutex) {
	defer wg.Done()
	lock.Lock() // Wait because someone is making a deposit	- After that when some goroutines use a balance variable wait
	b := balance
	// time.Sleep(time.Second * 5)
	balance = b + amount
	lock.Unlock() // unlock the lock goroutine when all go routines have been finished
}

func Balance(lock *sync.RWMutex) int {
	lock.RLock()
	b := balance
	lock.RUnlock()
	return b
}

// 1 Deposit() -> writing (Race condition)
// N Balance() -> read balance
func main() {
	var wg sync.WaitGroup
	var lock sync.RWMutex
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go Deposit(i*100, &wg, &lock)
	}
	fmt.Println(Balance(&lock))
	wg.Wait()
	fmt.Println(Balance(&lock))
}
