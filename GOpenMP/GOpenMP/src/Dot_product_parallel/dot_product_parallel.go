package main

import (
	"fmt"
)

import "runtime"

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func main() {
	_init_numCPUs()
	var sum float64
	var a [256]float64
	var b [256]float64
	var i int
	var n int = 256
	// Size estoy toqueteando comentarios
	for i = 0; i < n; i++ {
		a[i] = float64(i) * 0.5
		b[i] = float64(i) * 2.0
	}
	sum = 0
	var _barrier_0_float64 = make(chan float64)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			var sum float64
			for _i := _routine_num + 0; _i < (n+0)/1; _i += _numCPUs {
				sum += a[_i] * b[_i]
			}
			_barrier_0_float64 <- sum
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		sum += <-_barrier_0_float64
	}

	fmt.Println("a*b =", sum)
}
