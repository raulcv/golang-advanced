package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func Fibonacci(key int) int {
	if key <= 1 {
		return key
	}
	return Fibonacci(key-1) + Fibonacci(key-2)
}

// Memory holds a function and a map of results
type Memory struct {
	Funct Function               //Function to be used
	Cache map[int]FunctionResult //Map of results for a given key
	lock  sync.Mutex             //
}

// It function has to receive a value and return a value and an error
type Function func(key int) (interface{}, error)

// Mask for a result of a function
type FunctionResult struct {
	value interface{}
	err   error
}

// Mask for create a new cache
func NewCache(f Function) *Memory {
	return &Memory{
		Funct: f,
		Cache: make(map[int]FunctionResult),
	}
}

// It return the value for a given key value
func (m *Memory) GetCache(key int) (interface{}, error) {
	m.lock.Lock()
	result, exists := m.Cache[key] // chech if the key exists in cache already
	m.lock.Unlock()
	if !exists {
		m.lock.Lock()
		result.value, result.err = m.Funct(key) // Calculate the value for the key when It doesn't exist yet
		m.Cache[key] = result                   // Save the value in memory
		m.lock.Unlock()
	}
	return result.value, result.err
}

// Function calculate and return the value for a given key
func GetFibonacci(key int) (interface{}, error) {
	return Fibonacci(key), nil
}

func main() {
	//Create a cache and some values
	cache := NewCache(GetFibonacci)
	values := []int{42, 40, 41, 42, 38, 42, 42, 42}

	var wg sync.WaitGroup

	maxGoroutines := 2
	channel := make(chan int, maxGoroutines)
	for _, val := range values {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			channel <- 1
			start := time.Now()
			value, err := cache.GetCache(index)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("value: %d, Time: %s, Fibo Resul: %d\n", index, time.Since(start), value)

			<-channel

		}(val)
	}
	wg.Wait()
}
