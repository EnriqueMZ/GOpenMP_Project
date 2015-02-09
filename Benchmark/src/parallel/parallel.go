package main

import (
	"fmt"
	"runtime"
)

const Num_routines = 100

var ch = make(chan int)
var a int

func main() {
	runtime.GOMAXPROCS(4)
	a = 5.0
	for i := 0; i < Num_routines; i++ {
		go func(b int) {
			a = int(a) + 2
			fmt.Println("a: ", a, " desde la gorroutine ", b)
			ch <- 0
		}(i)
	}

	for i := 0; i < Num_routines; i++ {
		<-ch
	}

}
