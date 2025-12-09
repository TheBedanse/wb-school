package main

import (
	"fmt"
	"sync"
)

func main() {
	// Создаем два канала
	numbersChan := make(chan int)
	resChan := make(chan int)

	var wg sync.WaitGroup

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(numbersChan)

		for _, num := range numbers {
			numbersChan <- num
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(resChan)

		for num := range numbersChan {
			resChan <- num * 2
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for doubledNum := range resChan {
			fmt.Println(doubledNum)
		}
	}()

	wg.Wait()

}
