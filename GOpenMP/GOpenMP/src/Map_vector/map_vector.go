package main

import (
	"fmt"
	"gomp_lib"
	) 

func main() {
	var wtime float64
	var a [100000]float64
	var	b [100000]float64
	var	i int
	var	n int = 100000  // Size
	
	for i = 0; i < n; i++ {
		a[i] = float64(i) * 0.5
		b[i] = float64(i) * 2.0
		}
	wtime = gomp_lib.Gomp_get_wtime()
	//pragma gomp parallel for private(i) shared(b)
	for i = 0; i < n; i++ {
		 b[i] += a[i] * b[i]
	}
	wtime = gomp_lib.Gomp_get_wtime() - wtime
	fmt.Printf("Time: ")
	fmt.Printf("%14f\n", wtime)
}