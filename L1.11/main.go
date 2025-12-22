package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4}

	doubleNum(a, b)
	fmt.Println(doubleNum(a, b))

}

func doubleNum(a, b []int) []int {
	var res []int
	sudo := make(map[int]int)

	for _, k := range a {
		sudo[k]++
	}

	for _, k := range b {
		sudo[k]++
		if sudo[k] == 2 {
			res = append(res, k)
		}

	}
	return res
}
