/*
==================================================================
 SAXPY: Ejecución paralela mediante bloques de tamaño fijo
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
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	init := time.Now()
	flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            fmt.Println("Error: ", err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    
    if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            fmt.Println("Error: ", err)
        }
        pprof.WriteHeapProfile(f)
        f.Close()
        return
    }
    
	_init_numCPUs()
	var n int = 300000000 	// Tamaño de los vectores
	var a float32 = 2		// Factor de multiplicacion
	x, y := SAXPY_Init.Saxpy_init(n)
	
	init_p := time.Now()
	var _barrier_1_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num * (n / _numCPUs); i < (_routine_num + 1) * (n / _numCPUs); i++ { // Modo de paralelización.
				y[i] = a*x[i] + y[i]
			}
			_barrier_1_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_1_bool
	}
	fin_p := time.Since(init_p).Seconds()
	fin := time.Since(init).Seconds()
	fmt.Println(_numCPUs, ",", fin_p, ",", fin)

}