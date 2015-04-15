package main

import (
	"fmt"
	"gomp_lib"
)

import "runtime"

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func main() {
	_init_numCPUs()
	var wtime float64
	var a [100000]float64
	var b [100000]float64
	var i int
	var n int = 100000
	// Size estoy toqueteando comentarios
	for i = 0; i < n; i++ {
		a[i] = float64(i) * 0.5
		b[i] = float64(i) * 2.0
	}
	wtime = gomp_lib.Gomp_get_wtime()
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var (
				i int
			)
			for i = _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				var ()
				b[i] += a[i] * b[i]
			}
			_barrier_0_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_0_bool
	}

	wtime = gomp_lib.Gomp_get_wtime() - wtime
	fmt.Printf("Time: ")
	fmt.Printf("%14f\n", wtime)
}
