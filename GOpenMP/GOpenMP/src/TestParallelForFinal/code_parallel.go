package main

import (
	"fmt"
	//. "gomp_lib"
	"runtime"
)

func main() {
	_numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(_numCPUs)
	var sum1 int = 0
	var sum2 int = 0
	var prod float64 = 0
	var res float64 = 1000
	//var cont int = 0
	fmt.Println("Inicio de la region paralela")
	var _barrier_0_int = make(chan int)
	var _barrier_1_int = make(chan int)
	var _barrier_2_float64 = make(chan float64)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var (
				cont int
			)
			var sum1 int
			var sum2 int
			var prod float64
			for _i := _routine_num + 0; _i < (10+0)/2; _i += _numCPUs {
				sum1 += 1
				sum2 += 2
				prod *= 2
				res -= 2
				cont++
				fmt.Println("Gouroutine:", _routine_num, " cont =", cont)
				fmt.Println("Gouroutine del sistema:", runtime.NumGoroutine, " cont =", cont)
			}
			_barrier_0_int <- sum1
			_barrier_1_int <- sum2
			_barrier_2_float64 <- prod
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		sum1 += <-_barrier_0_int
		sum2 += <-_barrier_1_int
		prod *= <-_barrier_2_float64
	}

	fmt.Println("Fin de la region paralela")
	fmt.Println("Valores de la variable fuera del bloque parallel:", sum1, sum2, prod)
}
