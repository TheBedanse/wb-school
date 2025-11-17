package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go Worker(i, jobs, &wg, ctx)
	}

	go func() {
		defer close(jobs)
		counter := 1
		for {
			num := rand.Intn(1000)
			jobs <- num
			counter++
			time.Sleep(100 * time.Millisecond)
		}
	}()
	<-signalCh
	cancel()
	fmt.Println("Generate success")
	wg.Wait()
	fmt.Println("Finish workers")
}

func Worker(id int, jobs <-chan int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	for job := range jobs {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker id %d: %d\n", id, job)
			fmt.Printf("Worker id %d shuts down\n", id)
			return
		default:
			fmt.Printf("Worker id %d: %d\n", id, job)
		}
	}

}
