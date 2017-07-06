package lib

import (
	"math"
)

// output multiples of the provided number to the returned channel until the
// done channel is closed
func MultiplesOf(num int64, done chan struct{}) <-chan int64 {
	output := make(chan int64, 1)
	go func() {
		var cur int64
		for {
			cur += num
			select {
			case output <- cur:
			case <-done:
				close(output)
				return
			}
		}
	}()
	return output
}

// channel wrapper that provides the user with a Peek() method to view the next
// item available on the channel
type BufferedChannel struct {
	ok      bool
	next    int64
	channel <-chan int64
}

func NewBufferedChannel(channel <-chan int64) *BufferedChannel {
	next, ok := <-channel
	return &BufferedChannel{ok, next, channel}
}

func (bc *BufferedChannel) Peek() (next int64, ok bool) {
	next = bc.next
	ok = bc.ok
	return
}

func (bc *BufferedChannel) Receive() (next int64, ok bool) {
	next = bc.next
	ok = bc.ok
	bc.next, bc.ok = <-bc.channel
	return
}

// merge sorted input channels into a single output channel. if unique is true,
// duplicate values across channels will be squashed into a single value
func MergedSortedChannel(sortedChannels []chan int64, unique bool) <-chan int64 {
	output := make(chan int64)
	go func() {
		// get a list of BufferedChannels so we can peek at the next value
		chans := make([]*BufferedChannel, len(sortedChannels))
		for i, c := range sortedChannels {
			chans[i] = NewBufferedChannel(c)
		}

		lowestChans := make([]*BufferedChannel, 0, len(sortedChannels))
		for {
			var lowestValue int64
			// find the channels with the lowest next value
			for i, c := range chans {
				next, ok := c.Peek()
				if !ok {
					panic("here")
				} else {
					switch {
					case len(lowestChans) == 0:
						lowestValue = next
						fallthrough
					case next == lowestValue:
						lowestChans = append(lowestChans, chans[i])
					case next < lowestValue:
						lowestValue = next
						lowestChans[0] = c
						lowestChans = lowestChans[0:1]
					}
				}
			}

			// send the lowest value and receive on all channels that had
			// this value
			toSend := lowestValue
			for i, c := range lowestChans {
				if i == 0 || !unique {
					output <- toSend
				}
				c.Receive()
			}

			// clear out our list of lowest channels for reuse
			lowestChans = lowestChans[0:0]
		}
	}()
	return output
}

// output values for the fibonacci sequence on the returned channel
func Fib(done <-chan struct{}) <-chan int64 {
	output := make(chan int64)
	go func() {
		prev := int64(1)
		current := int64(1)
		output <- prev
		output <- current
		for {
			tmp := current
			current = prev + current
			prev = tmp
			select {
			case output <- current:
			case <-done:
				close(output)
				return
			}
		}
	}()
	return output
}

// return the number of digits in the nth fibonacci number
// NOTE: This function doesn't appear to work yet
func FibNDigits(term int) int {
	switch term {
	case 1:
		return 1
	default:
		return ((term - 2) / 5) + 1
	}
}

// output all prime numbers less than maxPrime on the returned channel
func Primes(maxPrime int64) <-chan int64 {
	output := make(chan int64)
	go func() {
		output <- 1
		output <- 2
		// list should be large enough to hold all odd numbers less than maxPrime
		numOdds := int64(math.Ceil(float64(maxPrime) / 2.0))
		isOddComposite := make([]bool, numOdds)
		for i := int64(3); i <= maxPrime; i += int64(2) {
			if isOddComposite[i/2] {
				// nothing to do because this number is not prime
			} else {
				// i is a prime number
				output <- i
				// mark all multiples of i as odd composites
				maxMultiple := maxPrime / i
				for j := int64(3); j <= maxMultiple; j += 2 {
					composite := i * j
					if composite > maxPrime {
						panic("bug")
					}
					isOddComposite[composite/2] = true
				}
			}
		}
		close(output)
	}()
	return output
}
