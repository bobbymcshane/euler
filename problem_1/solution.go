package main

import (
	"euler/lib"
	"fmt"
)

func main() {
	done := make(chan struct{})
	c := lib.MergedSortedChannel(
		[]chan int{lib.MultiplesOf(3, done), lib.MultiplesOf(5, done)}, true)

	var sum int
	for {
		next := <-c
		if next < 23 {
			sum += next
		} else {
			break
		}
	}
	fmt.Println(sum)
}
