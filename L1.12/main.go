package main

import "fmt"

func main() {
	sliceString := []string{"cat", "cat", "dog", "cat", "tree"}
	result := delDouble(sliceString)
	fmt.Println(result)

}

func delDouble(sliceString []string) []string {
	var res []string
	sudo := make(map[string]bool)
	for _, k := range sliceString {
		if !sudo[k] {
			sudo[k] = true
			res = append(res, k)
		}
	}

	return res

}
