package main

import (
	"fmt"
	) 

func main() {
	var sum float64
	var a [256]float64
	var	b [256]float64
	var	i int
	var	n int = 256  // Size
	
	for i = 0; i < n; i++ {
		a[i] = float64(i) * 0.5
		b[i] = float64(i) * 2.0
		}
	sum = 0
	
	//pragma gomp parallel for private(i) reduction(+:sum)
	for i = 0; i < n; i++ {
		sum += a[i] * b[i]
	}
	fmt. Println("a*b =", sum)
}