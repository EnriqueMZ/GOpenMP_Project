package Dot_Bench

import (
	//"runtime"
	//"fmt"
	)

func Dot_parallel_A(){
	_init_numCPUs()
	var sum float64 = 0
	n, a, b := Dot_Init_A() 
	var _barrier_0_float64 = make(chan float64)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var (
				i int
			)
			var sum float64
			for i = _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				sum += a[i] * b[i]
			}
			_barrier_0_float64 <- sum
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		sum += <-_barrier_0_float64
	}
	//fmt.Println("Parallel A result: ", sum)
}

func Dot_parallel_B(){
	_init_numCPUs()
	var sum float64 = 0
	var n int = 300000000
	a, b := Dot_Init_B(n) 
	var _barrier_0_float64 = make(chan float64)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var (
				i int
			)
			var sum float64
			for i = _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				sum += a[i] * b[i]
			}
			_barrier_0_float64 <- sum
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		sum += <-_barrier_0_float64
	}
	//fmt.Println("Parallel B result: ", sum)
}

