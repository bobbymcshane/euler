package main

import (
	"euler/lib"
	"fmt"
)

func main() {
	var i int
	for i = 1; lib.FibNDigits(i) < 1000; i++ {
	}
	fmt.Println(i)
}
