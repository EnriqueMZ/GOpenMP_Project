package main

import (
	"fmt"
	"gomp_lib"
	"math"
)

import "runtime"

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func main() {
	_init_numCPUs()
	var a [500][500]float64
	//var angle float64
	var b [500][500]float64
	var c [500][500]float64
	//var i int
	//var j int
	//var k int
	var n int = 500
	var pi float64 = 3.141592653589793
	var s float64
	var thread_num int
	var wtime float64
	fmt.Printf("\n")
	fmt.Printf("MXM_OPENMP:\n")
	fmt.Printf("  C/OpenMP version\n")
	fmt.Printf("  Compute matrix product C = A * B.\n")
	thread_num = gomp_lib.Gomp_get_num_routines()
	fmt.Printf("\n")
	fmt.Printf("  The number of processors available = %d\n", gomp_lib.Gomp_get_num_procs())
	fmt.Printf("  The number of threads available    = %d\n", thread_num)
	fmt.Printf("  The matrix order N                 = %d\n", n)
	s = 1.0 / math.Sqrt(float64(n))
	wtime = gomp_lib.Gomp_get_wtime()
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var (
				angle	float64
				//i	int
				j	int
				k	int
			)
			for _i := _routine_num + 0; _i < (n+0)/1; _i += _numCPUs {
				var ()
				for j = 0; j < n; j++ {
					angle = 2.0 * pi * float64(_i) * float64(j) / float64(n)
					a[_i][j] = s * (math.Sin(angle) + math.Cos(angle))
				}
			}
			for _i := _routine_num + 0; _i < (n+0)/1; _i += _numCPUs {
				var ()
				for j = 0; j < n; j++ {
					b[_i][j] = a[_i][j]
				}
			}
			for _i := _routine_num + 0; _i < (n+0)/1; _i += _numCPUs {
				var ()
				for j = 0; j < n; j++ {
					c[_i][j] = 0.0
					for k = 0; k < n; k++ {
						c[_i][j] = c[_i][j] + a[_i][k]*b[k][j]
					}
				}
			}
			_barrier_0_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_0_bool
	}

	wtime = gomp_lib.Gomp_get_wtime() - wtime
	fmt.Printf("  Elapsed seconds = %g\n", wtime)
	fmt.Printf("  C(100,100)  = %g\n", c[99][99])
	fmt.Printf("\n")
	fmt.Printf("MXM_OPENMP:\n")
	fmt.Printf("  Normal end of execution.\n")
	fmt.Printf("\n")
}
