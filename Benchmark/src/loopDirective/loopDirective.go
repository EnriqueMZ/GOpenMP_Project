package main

import (
	"fmt"
	. "gomp_lib"
	"runtime"
)

const num_steps = 10000
var problema int
var ch = make(chan int)
func main() {
	fmt.Println(runtime.NumGoroutine(),runtime.NumCPU());
	runtime.GOMAXPROCS(4)
	//pragma gomp parallel num_routines(4)
	Gomp_set_num_routines(4)
	g:=Gomp_get_num_routines() //usamos una variable y le damos el valor de la llamada al metodo de la libreria
	for i := 0; i < Gomp_get_num_routines(); i++ {
		go func(j int) {
			tid := j
			//pragma gomp for
			for i := tid * (num_steps / g); i < (num_steps/g)*(tid+1); i++ {
				a:=99
				a++
				problema++
			}
			ch <- 0
		}(i)
	}
	for i := 0; i < Gomp_get_num_routines(); i++ {
		<-ch
	}
	fmt.Println("problema: ", problema, " expected: ", num_steps)

}
