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

func Dot_Init_B_simple(size int) ([]float64, []float64) {
	
	a := make([]float64, size)
	b := make([]float64, size)

	for i := 0; i < size; i++ {
		a[i] = 1
		b[i] = 1
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
			for i = _routine_num * (n / _numCPUs); i < (_routine_num+1)*(n/_numCPUs); i++ {
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
	var n int = 24000000
	a, b := Dot_Init_B_simple(n) 
	init_p = time.Now()
	var _barrier_0_float64 = make(chan float64)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int, _a []float64, _b []float64, __numCPUs int, __barrier_0_float64 chan float64) {
			var (
				i int
				cont int = 0
			)
			var sum float64
			for i = 0; i < len(_a); i++ {
				sum += _a[i] * _b[i]
			}
			cont += 
			__barrier_0_float64 <- sum
		}(_i, a[(_i*8)+cont:(_i+1)*8], b[_i*8:(_i+1)*8], _numCPUs, _barrier_0_float64)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		sum += <-_barrier_0_float64
	}
	fin_p := time.Since(init_p).Seconds()
	fin := time.Since(ini).Seconds()
	fmt.Println(_numCPUs, ",", fin_p, ",", fin, "," , sum)
}

func main() {
	//Dot_parallel_A()
	Dot_parallel_B()
}
