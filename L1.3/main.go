package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Use go run main.go (number_workers)")
		return
	}

	numWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil || numWorkers <= 0 {
		fmt.Printf("Error number of workers: %v", err)
		return
	}

	jobs := make(chan int)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go Worker(i, jobs, &wg)
	}

	go func() {
		defer close(jobs)
		counter := 1
		for {
			num := rand.Intn(1000)
			jobs <- num
			counter++
			if counter > 100 {
				break
			}
		}
	}()
	fmt.Println("Generate success")
	wg.Wait()
	fmt.Println("Finish workers")
}

func Worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker id %d: %d\n", id, job)
	}

}
