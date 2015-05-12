package main

import (
	"runtime"
	"fmt"
	"time"
	)

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}

func Dot_Init_A() (int, [10000000]float64, [10000000]float64) {
	var a [10000000]float64
	var	b [10000000]float64
	var	n int = 10000000  // Size
	
	for i := 0; i < n; i++ {
		a[i] = float64(i) * 0.5
		b[i] = float64(i) * 2.0
		}
	
	return n, a, b
}

func Dot_Init_B(size int) ([]float64, []float64) {
	
	a := make([]float64, size)
	b := make([]float64, size)

	for i := 0; i < size; i++ {
		a[i] = float64(i) * 0.5
		b[i] = float64(i) * 2.0
	}
	return a, b
}

func Dot_parallel_A(){
	_init_numCPUs()
	var sum float64 = 0
	n, a, b := Dot_Init_A()
	init := time.Now() 
	var _barrier_0_float64 = make(chan float64)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var (
				i int
			)
			var sum float64
			for i = _routine_num * (n / _numCPUs); i < (_routine_num+1)*(n/_numCPUs); i++ {
				sum += a[i] * b[i]
			}
			_barrier_0_float64 <- sum
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		sum += <-_barrier_0_float64
	}
	fin := time.Since(init)
	fmt.Println("Time: ", fin)
	fmt.Println("Parallel A result: ", sum)
}

func Dot_parallel_B(){
	_init_numCPUs()
	var sum float64 = 0
	var n int = 300000000
	a, b := Dot_Init_B(n)
	init := time.Now() 
	var _barrier_0_float64 = make(chan float64)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var (
				i int
			)
			var sum float64
			for i = _routine_num * (n / _numCPUs); i < (_routine_num+1)*(n/_numCPUs); i++ {
				sum += a[i] * b[i]
			}
			_barrier_0_float64 <- sum
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		sum += <-_barrier_0_float64
	}
	fin := time.Since(init)
	fmt.Println("Time: ", fin)
	fmt.Println("Parallel B result: ", sum)
}

func main() {
	//Dot_parallel_A()
	Dot_parallel_B()
}
