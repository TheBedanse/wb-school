package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	N := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), N)
	defer cancel()
	var wg sync.WaitGroup
	chanel := make(chan string)

	wg.Add(1)
	go Write(ctx, chanel, &wg)

	wg.Add(1)
	go Read(ctx, chanel, &wg)

	wg.Wait()
	fmt.Println("Programm finish")
}

func Write(ctx context.Context, write chan<- string, wg *sync.WaitGroup) {
	defer close(write)
	defer wg.Done()
	for {
		generateTime := time.Now().String()
		select {
		case <-ctx.Done():
			fmt.Println("3 Second later, close write")
			return
		case write <- generateTime:
			fmt.Printf("Write in chan: %s\n", generateTime)
			time.Sleep(50 * time.Millisecond)
		}
	}

}

func Read(ctx context.Context, read <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("3 Second later, close read")
			return
		case t := <-read:
			fmt.Printf("read %v\n", t)
		}
	}
}
