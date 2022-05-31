package main

import (
	"fmt"
	"sync"
	"time"
)

// New way to save a information in Cache for Fibonacci number calculation || Reutilisation de computacion intensiva
// Receiving information to save in memory
func ExpensiveFinoacci(n int) int {
	fmt.Printf("Calculating Expensive Fibonacci of %d\n", n)
	time.Sleep(time.Second * 5) // Emulate expensive calculation of Fibonacci
	return n                    // return
}

type Service struct {
	InProgress map[int]bool       // Mapping of int keys to in progress and bool to indicate if the key is in progress
	IsPending  map[int][]chan int // Mapping of int keys to slice of channels to notify when the key is ready
	Lock       sync.RWMutex       // Protects the in progress map – lock race condition
}

func (s *Service) Work(job int) {
	s.Lock.RLock()              // Locking read services
	exists := s.InProgress[job] // Check if the job is already in progress
	if exists {
		s.Lock.RUnlock()           // Unlock the in progress if it is not in progress
		response := make(chan int) // Create a channel to wait for the response from the worker
		defer close(response)      // Close the channel when the function is done

		s.Lock.Lock()                                         // Lock the read and write service to add the channel to the pending map
		s.IsPending[job] = append(s.IsPending[job], response) // Add the channel to the pending map
		s.Lock.Unlock()                                       // Unlock the read and write service
		fmt.Printf("Job %d is already in progress\n", job)
		resp := <-response                             // Wait for the response from the worker
		fmt.Printf("Job %d is done, received\n", resp) // Print the response
		return                                         // Return the response
	}
	s.Lock.RUnlock() // Unlock the read service

	s.Lock.Lock()            // Lock the read and write service to add the job to the in progress map
	s.InProgress[job] = true // Add the job to the in progress map, some other worker did it before us
	s.Lock.Unlock()          // Unlock the read and write service

	fmt.Printf("Job %d is in progress\n", job) // Print the job is in progress
	result := ExpensiveFinoacci(job)           // Calculate the result

	s.Lock.RLock()                             // Lock the read service
	pendingWorkers, exists := s.IsPending[job] // Check if the job is still in progress
	s.Lock.RUnlock()                           // Unlock the read service
	if exists {
		for _, pendingWorker := range pendingWorkers {
			pendingWorker <- result // Notify the pending workers
		}
		fmt.Printf("Result sent – all pending workers notified | job: %d\n", job)
	}
	s.Lock.Lock()                          // Lock the read and write service
	s.InProgress[job] = false              // set the job to not in progress, It have been calculated
	s.IsPending[job] = make([]chan int, 0) // Remove the pending workers || empty the pending workers
	s.Lock.Unlock()                        // Unlock the read and write service
}

func NewService() *Service {
	return &Service{
		InProgress: make(map[int]bool),
		IsPending:  make(map[int][]chan int),
	}
}
func main() {
	service := NewService()
	jobs := []int{3, 4, 5, 5, 4, 8, 8, 8}
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, n := range jobs {
		go func(job int) {
			defer wg.Done()
			service.Work(job)
		}(n)
	}
	wg.Wait()
}
