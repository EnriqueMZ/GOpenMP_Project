package main

import (
	"fmt"
	"time"
	"math"
)

import "runtime"

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func main() {
	init := time.Now()
	_init_numCPUs()
	var a [1000][1000]float64
	var angle float64
	var b [1000][1000]float64
	var c [1000][1000]float64
	var i int
	var j int
	var k int
	var n int = 1000
	var pi float64 = 3.141592653589793
	var s float64
	s = 1.0 / math.Sqrt(float64(n))
	
	init_p := time.Now()
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int, _angle float64, _i int, _j int, _k int) {
			var (
				angle	float64
				i	int
				j	int
				k	int
			)
			for i = _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				var ()
				for j = 0; j < n; j++ {
					angle = 2.0 * pi * float64(i) * float64(j) / float64(n)
					a[i][j] = s * (math.Sin(angle) + math.Cos(angle))
				}
			}
			for i = _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				var ()
				for j = 0; j < n; j++ {
					b[i][j] = a[i][j]
				}
			}
			for i = _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				var ()
				for j = 0; j < n; j++ {
					c[i][j] = 0.0
					for k = 0; k < n; k++ {
						c[i][j] = c[i][j] + a[i][k]*b[k][j]
					}
				}
			}
			_barrier_0_bool <- true
		}(_i, angle, i, j, k)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_0_bool
	}
	fin_p := time.Since(init_p).Seconds()
	fin := time.Since(init).Seconds()
	fmt.Println(_numCPUs, ",", fin_p, ",", fin)
}
