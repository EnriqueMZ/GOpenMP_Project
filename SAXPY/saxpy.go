// Alberto Castaño

package main

import (
	"os"
	"strconv"
	"math/rand"
)

import "runtime"

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func main() {
	_init_numCPUs()
	var n int
	n, _ = strconv.Atoi(os.Args[1])
	var a float32 = 2
	x := make([]float32, n)
	y := make([]float32, n)
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				x[i] = float32(rand.Int())
				y[i] = float32(rand.Int())
			}
			_barrier_0_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_0_bool
	}

	var _barrier_1_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				y[i] = a*x[i] + y[i]
			}
			_barrier_1_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_1_bool
	}

}
