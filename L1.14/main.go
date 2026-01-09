package main

import (
	"fmt"
	"reflect"
)

func main() {
	ch := make(chan int)
	checkType(10)
	checkType("Hi")
	checkType(true)
	checkType(ch)

}

func checkType(x any) {
	switch v := x.(type) {
	case int:
		fmt.Printf("This int: %d\n", v)
	case string:
		fmt.Printf("This string: %s\n", v)
	case bool:
		fmt.Printf("This bool: %v\n", v)
	default:
		t := reflect.TypeOf(x)
		if t != nil && t.Kind() == reflect.Chan {
			fmt.Printf("This chan: %T\n", v)
		} else {
			fmt.Println("Unknown type")
		}

	}
}
