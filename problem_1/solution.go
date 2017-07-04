package main

import (
	"euler/lib"
	"fmt"
)

func main() {
	done := make(chan struct{})
	c := lib.MergedSortedChannel(
		[]chan int64{lib.MultiplesOf(3, done), lib.MultiplesOf(5, done)}, true)

	var sum int64
	for {
		next, ok := <-c
		if ok {
			if next < 1000 {
				sum += next
			} else {
				close(done)
				break
			}
		}
	}
	fmt.Println(sum)
}
