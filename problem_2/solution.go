package main

import (
	"euler/lib"
	"fmt"
)

func main() {
	done := make(chan struct{})
	fib := lib.Fib(done)
	var sum int
	for {
		term := <-fib
		if term > 4000000 {
			break
		} else if term%2 == 0 {
			sum += term
		}
	}
	fmt.Println(sum)
}
