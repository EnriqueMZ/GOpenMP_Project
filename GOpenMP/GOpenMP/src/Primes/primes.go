package main

import (
	"fmt"
	"gomp_lib"
)

func prime_number(n int) int {
	var i int
	var j int
	var prime int
	var total int = 0

	//pragma gomp parallel for shared(n) private(i, j, prime) reduction (+:total)
		for i = 2; i <= n; i++ {
			prime = 1
			for j = 2; j < i; j++ {
				if i%j == 0 {
					prime = 0
					break
				}
			}
			total = total + prime
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

	var	n_factor int
	var	n_hi int
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

	prime_number_sweep ( n_lo, n_hi, n_factor )

	n_lo = 5
	n_hi = 500000
	n_factor = 10

	prime_number_sweep ( n_lo, n_hi, n_factor )

	fmt.Print("\n")
	fmt.Println("PRIME_OPENMP")
	fmt.Println("  Normal end of execution.")
}
