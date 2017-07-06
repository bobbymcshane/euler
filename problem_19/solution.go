package main

import (
	"fmt"
	"time"
)

// strategy 1: brute force. figure out the day on the first of every month
func isSunday(t time.Time) bool {
	return t.Weekday() == time.Sunday
}

func main() {
	var numSundays int
	// iterate over the first of each month and count the number of Sundays we find
	for y := 1901; y <= 2000; y++ {
		for m := 1; m <= 12; m++ {
			if isSunday(time.Date(y, time.Month(m), 1, 1, 0, 0, 0, time.UTC)) {
				numSundays++
			}
		}
	}
	fmt.Println(numSundays)
}
