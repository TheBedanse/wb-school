package main

import (
	"Gorutines/exit"
	"fmt"
)

func main() {
	fmt.Println("Demo stop gorutines:")

	fmt.Println("==========ByChan==============")
	exit.ByChan()

	fmt.Println("==========ByContext===========")
	exit.ByContext()

	fmt.Println("==========ByChanClose=========")
	exit.ByChanClose()

	fmt.Println("==========ByCondition=========")
	exit.ByCondition()

	fmt.Println("==========ByGoexit============")
	exit.ByGoexit()

	fmt.Println("==========ByPanic=============")
	exit.ByPanic()
}
