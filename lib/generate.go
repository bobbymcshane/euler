package lib

// output multiples of the provided number to the returned channel until the
// done channel is closed
func MultiplesOf(num int, done chan struct{}) <-chan int {
	output := make(chan int, 1)
	go func() {
		var cur int
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
	next    int
	channel <-chan int
}

func NewBufferedChannel(channel <-chan int) *BufferedChannel {
	next, ok := <-channel
	return &BufferedChannel{ok, next, channel}
}

func (bc *BufferedChannel) Peek() (next int, ok bool) {
	next = bc.next
	ok = bc.ok
	return
}

func (bc *BufferedChannel) Receive() (next int, ok bool) {
	next = bc.next
	ok = bc.ok
	bc.next, bc.ok = <-bc.channel
	return
}

// merge sorted input channels into a single output channel. if unique is true,
// duplicate values across channels will be squashed into a single value
func MergedSortedChannel(sortedChannels []chan int, unique bool) <-chan int {
	output := make(chan int)
	go func() {
		// get a list of BufferedChannels so we can peek at the next value
		chans := make([]*BufferedChannel, len(sortedChannels))
		for i, c := range sortedChannels {
			chans[i] = NewBufferedChannel(c)
		}

		lowestChans := make([]*BufferedChannel, 0, len(sortedChannels))
		for {
			var lowestValue int
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
func Fib(done <-chan struct{}) <-chan int {
	output := make(chan int)
	go func() {
		prev := 1
		current := 2
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
