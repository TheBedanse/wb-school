package exit

import (
	"fmt"
	"time"
)

func ByCondition() {
	go condition()

	time.Sleep(10 * time.Millisecond)
}

func condition() {
	for i := 0; i < 3; i++ {
		fmt.Printf("Gorutine cycle: %d\n", i)
	}
	fmt.Println("Goruntine exit")
}

// Выход Горутины по условию
