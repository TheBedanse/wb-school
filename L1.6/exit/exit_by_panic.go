package exit

import (
	"fmt"
	"time"
)

func ByPanic() {
	go panicRec()
	time.Sleep(50 * time.Millisecond)
}

func panicRec() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recover panic:%v\n", r)
			fmt.Println("Gorutine exit")
		}
	}()
	for i := 0; i < 4; i++ {
		fmt.Printf("Gorutine work %d\n", i)
		if i == 2 {
			panic("Error panic")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

//Выход горутины из-за паники
