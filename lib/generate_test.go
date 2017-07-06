package lib

import (
	"testing"
)

func TestPrimes(t *testing.T) {
	primes := Primes(int64(100))
	primesUnder100 := []int64{1, 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	for _, p := range primesUnder100 {
		prime, ok := <-primes
		if !ok {
			t.Fatalf("channel closed prematurely. stopped at %v", p)
		}

		if prime != p {
			t.Fatalf("mismatched prime number; expected %v, got %v", p, prime)
		}
	}

	p, ok := <-primes
	if ok {
		t.Fatalf("channel remained open after it should have closed. got %v", p)
	}
}
