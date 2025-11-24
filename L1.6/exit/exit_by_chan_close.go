package exit

import (
	"fmt"
	"time"
)

func ByChanClose() {
	data := make(chan int)
	go chanClose(data)
	for i := 0; i < 3; i++ {
		data <- i
		time.Sleep(10 * time.Millisecond)
	}
	close(data)
	time.Sleep(10 * time.Millisecond)

}

func chanClose(data <-chan int) {
	for {
		data, ok := <-data
		if !ok {
			fmt.Println("Channel close, exit gorutine")
			return
		}
		fmt.Printf("Gorunine working, data:%d\n", data)
	}
}

//Выход горутины через закрытие канала
