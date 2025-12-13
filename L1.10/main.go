package main

import "fmt"

func main() {
	temp := [8]float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	var groupMinusTwenty []float64
	var groupTen []float64
	var groupTwenty []float64
	var groupThirty []float64

	for _, i := range temp {
		if i <= -20 && i > -30 {
			groupMinusTwenty = append(groupMinusTwenty, i)
		} else if i < 20 && i >= 10 {
			groupTen = append(groupTen, i)
		} else if i < 30 && i >= 20 {
			groupTwenty = append(groupTwenty, i)
		} else if i < 40 && i >= 30 {
			groupThirty = append(groupThirty, i)
		}
	}
	fmt.Printf("-20:%v\n", groupMinusTwenty)
	fmt.Printf("10:%v\n", groupTen)
	fmt.Printf("20:%v\n", groupTwenty)
	fmt.Printf("30:%v\n", groupThirty)
}
