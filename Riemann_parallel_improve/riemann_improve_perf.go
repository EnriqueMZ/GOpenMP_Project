
package main

import (
	"os"
	"strconv"
	"fmt"
	"time"
)

import "runtime"

var step float64
var num_steps int
var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func main() {
	init := time.Now()
	_init_numCPUs()
	num_steps, _ = strconv.Atoi(os.Args[1])
	var x, sum float64
	step = 1.0 / float64(num_steps)
	init_p := time.Now()
	var _barrier_0_float64 = make(chan float64)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			var sum float64
			for i := _routine_num * (num_steps / _numCPUs); i < (_routine_num+1)*(num_steps/_numCPUs); i++ {
				x = (float64(i) + 0.5) * step
				sum = sum + 4.0/(1.0+x*x)
			}
			_barrier_0_float64 <- sum
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		sum += <-_barrier_0_float64
	}
	sum = step * sum
	//fmt.Println(sum)
	fin_p := time.Since(init_p).Seconds()
	fin := time.Since(init).Seconds()
	fmt.Println(_numCPUs, ",", fin_p, ",", fin)
}
