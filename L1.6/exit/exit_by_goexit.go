package exit

import (
	"fmt"
	"runtime"
	"time"
)

func ByGoexit() {
	go goexit()
	time.Sleep(50 * time.Millisecond)

}

func goexit() {
	defer fmt.Println("Gorutine exit")
	for i := 0; i < 4; i++ {
		fmt.Printf("Gorutine work %d\n", i)
		if i == 2 {
			time.Sleep(10 * time.Millisecond)
			runtime.Goexit()
		}
	}
}

//Выход горутины с помощью runtime.Goexit
