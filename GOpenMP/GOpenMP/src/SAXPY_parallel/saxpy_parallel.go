package main

import (
	"fmt"
)

import "runtime"

func main() {
	_numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(_numCPUs)
	var n int = 10
	var a float64 = 2
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	y := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("Vector x antes del parallel:", x)
	fmt.Println("Vector y antes del parallel:", y)
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for _i := _routine_num; _i < ((n+0)-0)/1; _i += _numCPUs {
				y[_i] = a*x[_i] + y[_i]
			}
			_barrier_0_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_0_bool
	}

	fmt.Println("Vector x despues del parallel:", x)
	fmt.Println("Vector y despues del parallel:", y)
}
