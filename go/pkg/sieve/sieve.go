// Working concurrency segmentted sieve
// Be sure to be in cd \go
// run "go test ./..." in terminal
// run "go test -count=1 ./..." to disable test caching and force the rerun with timed return
// Resource: https://cp-algorithms.com/algebra/sieve-of-eratosthenes.html
package sieve

import (
	"math"
	"runtime"
	"sync"
)

type Sieve struct{}

func NewSieve() *Sieve {
	return &Sieve{}
}

func (s *Sieve) NthPrime(n int64) int64 { // API call which determines which version of sieve to run
	if n < 0 {
		n = 0
	}
	if n >= 10_000_000 {
		return s.segmentedConcurrentSieve(n)
	}
	return s.simpleSieve(n)
}

func (s *Sieve) simpleSieve(n int64) int64 { // Standard implementation of Sieve of Eratosthenes
	var uBound int
	if n < 10 {
		uBound = 30
	} else {
		m := float64(n)
		uBound = int(m*(math.Log(m)+math.Log(math.Log(m))) + 2)
	}

	isPrime := make([]bool, uBound+1)
	for i := 2; i <= uBound; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= uBound; i++ {
		if isPrime[i] {
			for j := i * i; j <= uBound; j += i {
				isPrime[j] = false
			}
		}
	}

	count := int64(0)
	for i := 2; i <= uBound; i++ {
		if isPrime[i] {
			if count == n {
				return int64(i)
			}
			count++
		}
	}
	return -1
}

func (s *Sieve) segmentedConcurrentSieve(n int64) int64 { // Concurrent implemenation with segmented approach
	m := float64(n)
	limit := int(m*(math.Log(m)+math.Log(math.Log(m))) + 2)
	segmentSize := 1_000_000

	sqrtLimit := int(math.Sqrt(float64(limit)))
	isPrimeSmall := make([]bool, sqrtLimit+1)
	for i := 2; i <= sqrtLimit; i++ {
		isPrimeSmall[i] = true
	}
	for i := 2; i*i <= sqrtLimit; i++ {
		if isPrimeSmall[i] {
			for j := i * i; j <= sqrtLimit; j += i {
				isPrimeSmall[j] = false
			}
		}
	}

	basePrimes := []int{}
	for i := 2; i <= sqrtLimit; i++ {
		if isPrimeSmall[i] {
			basePrimes = append(basePrimes, i)
		}
	}

	count := int64(0)
	for low := 2; low <= limit; low += segmentSize {
		high := min(low+segmentSize-1, limit)
		mark := make([]bool, high-low+1)
		for i := range mark {
			mark[i] = true
		}

		numThreads := runtime.NumCPU()
		var wg sync.WaitGroup // Synchronization
		wg.Add(numThreads)

		// Parallelism implementation
		// No need for mutex as segments do not overlap
		chunkSize := (high - low + 1) / numThreads
		for t := 0; t < numThreads; t++ { // Launch go routine per segmented range of 1 million
			start := low + t*chunkSize
			end := start + chunkSize - 1
			if t == numThreads-1 {
				end = high
			}
			go func(start, end int) { // Each go routine is responsible for its specified segment
				defer wg.Done()
				for _, p := range basePrimes {
					first := max(p*p, ((start+p-1)/p)*p)
					for j := first; j <= end; j += p {
						mark[j-low] = false
					}
				}
			}(start, end)
		}

		wg.Wait()

		for i := low; i <= high; i++ {
			if mark[i-low] {
				if count == n {
					return int64(i)
				}
				count++
			}
		}
	}
	return -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
