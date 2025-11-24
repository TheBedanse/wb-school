package exit

import (
	"context"
	"fmt"
	"time"
)

func ByContext() {
	ctx, cancel := context.WithCancel(context.Background())

	go conext(ctx)

	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(5 * time.Millisecond)
}

func conext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Gorutine received a signal stop")
			fmt.Println("Gorutine exit")
			return
		default:
			fmt.Println("Gorutine working")
			time.Sleep(10 * time.Millisecond)
		}
	}
}

//Выход горутины через context, так же можно использовать WithTimeout, WithDeadline
