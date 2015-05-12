package main

import (
	"flag"
    "fmt"
    "runtime/pprof"
    "os"
    "SAXPY_Init"
    "time"
    )

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
    
	var n int = 300000000
	var a float64 = 2
	
	x, y := SAXPY_Init.Saxpy_init(n)
	
	init := time.Now()
	for i:= 0; i < n; i++ {
		y[i] = a * x[i] + y[i]
		}
	fin := time.Since(init)
	fmt.Println("Time: ", fin)
}