package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// struct used to distribute the work of scoring a name at a given position in
// a sorted list
type namePosition struct {
	name     string
	position int
}

func (n *namePosition) Value() int {
	var sum int
	for _, c := range []byte(n.name) {
		sum += int(c - 'a' + 1)
	}
	return sum * n.position
}

// read from nameChan, compute the value of the namePosition, and pass the result to valueChan
func NewWorker(nameChan <-chan namePosition, valueChan chan int) {
	go func() {
		var sum int
		for {
			select {
			case namePos, ok := <-nameChan:
				if !ok {
					valueChan <- sum
					return
				} else {
					sum += namePos.Value()
				}
			}
		}
	}()
}

func main() {
	f, err := os.Open("./testdata/names.txt")
	if err != nil {
		panic(err)
	}

	var names sort.StringSlice

	// create a new scanner and read the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}

	// check for errors
	if err = scanner.Err(); err != nil {
		panic(err)
	}

	// close our file as we no longer need it
	f.Close()

	// sort our list of strings alphabetically
	names.Sort()

	// dispatch each of our names to worker goroutines to calculate their value
	numWorkers := 4
	nameChan := make(chan namePosition)
	valueChan := make(chan int)

	for i := 0; i < numWorkers; i++ {
		NewWorker(nameChan, valueChan)
	}

	for i, name := range names {
		nameChan <- namePosition{strings.ToLower(name), i + 1}
	}
	close(nameChan)

	var total int
	for i := 0; i < numWorkers; i++ {
		total += <-valueChan
	}

	fmt.Println(total)
}
