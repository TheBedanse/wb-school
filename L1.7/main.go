package main

import (
	"fmt"
	"sync"
)

func main() {
	var mapa sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value := fmt.Sprintf("Hi %d", i)
			mapa.Store(i, value)
		}()
	}
	wg.Wait()

	fmt.Println("Range map:")

	mapa.Range(func(key, value any) bool {
		fmt.Printf("key: %d, value: %s\n", key, value)
		return true
	})
}
