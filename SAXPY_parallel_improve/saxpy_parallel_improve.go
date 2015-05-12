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
	var n int = 300000000
	var a float64 = 2
	x, y := SAXPY_Init.Saxpy_init(n)
	
	init := time.Now()
	var _barrier_1_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num * (n / _numCPUs); i < (_routine_num+1)*(n/_numCPUs); i++ {
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