/*
==================================================================
 SAXPY: Ejecución paralela mediante distribucion modular
==================================================================
*/

package main

import (
	"runtime"
	"flag"
    "fmt"
    "runtime/pprof"
    "os"
    "SAXPY_Init"
    "time"
    )

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            fmt.Println("Error: ", err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
	_init_numCPUs()
	var n int = 300000000 	// Tamaño de los vectores
	var a float64 = 2		// Factor de multiplicacion
	x, y := SAXPY_Init.Saxpy_init(n)
	
	init := time.Now()
	var _barrier_1_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num + 0; i < (n+0)/1; i += _numCPUs { // Modo de paralelización.
				y[i] = a*x[i] + y[i]
			}
			_barrier_1_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_1_bool
	}
	fin := time.Since(init)
	fmt.Println("Time: ", fin)

}