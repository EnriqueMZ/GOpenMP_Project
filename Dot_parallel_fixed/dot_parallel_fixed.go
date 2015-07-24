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

var ini time.Time
var init_p time.Time
var fin time.Time
var fin_p time.Time

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
	ini = time.Now()
	_init_numCPUs()
	var sum float64 = 0
	n, a, b := Dot_Init_A()
	init_p = time.Now()
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
	fin_p := time.Since(init_p).Seconds()
	fin := time.Since(ini).Seconds()
	fmt.Println(_numCPUs, ",", fin_p, ",", fin)
}

func Dot_parallel_B(){
	ini = time.Now()
	_init_numCPUs()
	var sum float64 = 0
	var n int = 300000000
	a, b := Dot_Init_B(n) 
	init_p = time.Now()
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
	fin_p := time.Since(init_p).Seconds()
	fin := time.Since(ini).Seconds()
	fmt.Println(_numCPUs, ",", fin_p, ",", fin)
}

func main() {
	//Dot_parallel_A()
	Dot_parallel_B()
}
