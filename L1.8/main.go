package main

import (
	"fmt"
	"strconv"
)

func main() {

	num := int64(50)  // num
	bitNum := uint(6) //bit num
	bitVal := 0       // change bit value 0 or 1

	i := bitNum - 1

	result := SetBit(num, i, bitVal)

	fmt.Printf("Original num: %d (%s)\n",
		num, strconv.FormatInt(num, 2))
	fmt.Printf("Install %d-bit change to %d\n", bitNum, bitVal)
	fmt.Printf("Result: %d (%s)\n",
		result, strconv.FormatInt(result, 2))
}

func SetBit(n int64, i uint, bitVal int) int64 {
	if bitVal == 1 {
		return n | (1 << i)
	}
	return n &^ (1 << i)
}
