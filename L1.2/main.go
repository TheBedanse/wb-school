package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	nums := [5]int{2, 4, 6, 8, 10}

	for _, num := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("%d^2 = %d\n", num, Square(num))
		}()
	}
	wg.Wait()
}

func Square(a int) int {
	return a * a
}
