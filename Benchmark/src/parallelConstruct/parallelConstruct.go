package main

import (
	"fmt"
	. "gomp_lib"
	"runtime"
//	"sync"
)

var a int
var b int
var c int
//var mutex=sync.Mutex {}
func main() {
	runtime.GOMAXPROCS(4)
	ch := make(chan int)
	Gomp_set_num_routines(10000)
	g:=Gomp_get_num_routines()
	for i := 0; i < g; i++ {
		go func(j int) {
			tid := j
			a := 0
			b := 0
			a++
			b++
//			mutex.Lock()
			c++
//			mutex.Unlock()
			fmt.Println("El valor de las variables en la rutina", tid, "es: ", a, b, c)
			ch <- 0
		}(i)
	}
	for i := 0; i < g; i++ {
		<-ch
	}
	fmt.Println("El valor de las variables es: ", a, b, c)
}
