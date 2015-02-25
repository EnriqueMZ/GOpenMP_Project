package main

import (
	"fmt"
	"gomp_lib"
)

import "runtime"

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func prime_number(a int) int {
	var n int = a
	// var i int
	// var j int
	// var prime int
	var total int = 0
	var _barrier_0_int = make(chan int)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var (
				prime	int
				j	int
			)
			var total int
			for _i := _routine_num + 2; _i < (n+1)/1; _i += _numCPUs {
				prime = 1
				for j = 2; j < _i; j++ {
					if _i%j == 0 {
						prime = 0
						break
					}
				}
				total = total + prime
			}
			_barrier_0_int <- total
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		total += <-_barrier_0_int
	}

	return total
}
func prime_number_sweep(n_lo int, n_hi int, n_factor int) {
	var n int
	var primes int
	var wtime float64
	fmt.Print("\n")
	fmt.Println("TEST01")
	fmt.Println("  Call PRIME_NUMBER to count the primes from 1 to N.")
	fmt.Print("\n")
	fmt.Println("         N        Pi        Time")
	fmt.Print("\n")
	n = n_lo
	for n <= n_hi {
		wtime = gomp_lib.Gomp_get_wtime()
		primes = prime_number(n)
		wtime = gomp_lib.Gomp_get_wtime() - wtime
		fmt.Printf("  %8d %8d %14f \n", n, primes, wtime)
		n = n * n_factor
	}
}
func main() {
	_init_numCPUs()
	var n_factor int
	var n_hi int
	var n_lo int
	fmt.Print("\n")
	fmt.Println("PRIME_GOPENMP")
	fmt.Println("  Go/OpenMP version")
	fmt.Print("\n")
	fmt.Println("  Number of processors available = ", gomp_lib.Gomp_get_num_procs())
	fmt.Println("  Number of threads =              ", gomp_lib.Gomp_get_num_routines())
	n_lo = 1
	n_hi = 131072
	n_factor = 2
	prime_number_sweep(n_lo, n_hi, n_factor)
	n_lo = 5
	n_hi = 500000
	n_factor = 10
	prime_number_sweep(n_lo, n_hi, n_factor)
	fmt.Print("\n")
	fmt.Println("PRIME_OPENMP")
	fmt.Println("  Normal end of execution.")
}
