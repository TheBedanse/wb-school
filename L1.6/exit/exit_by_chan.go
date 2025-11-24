package exit

import (
	"fmt"
	"time"
)

func ByChan() {
	stop := make(chan bool)
	go channel(stop)

	time.Sleep(50 * time.Millisecond)
	stop <- true
}

func channel(stop <-chan bool) {
	for {
		select {
		case <-stop:
			fmt.Println("Gorutine received a signal stop")
			fmt.Println("Gorutine exit")
			return
		default:
			fmt.Println("Gorutine working")
			time.Sleep(20 * time.Millisecond)
		}
	}

}

//Выход горутины через канал(сигнал)
