package main

import (
	"fmt"
	"runtime"
	."gomp_lib"
)

var n=100
var wch = make(chan int)
var sum =5
func main() {
//pragma gomp parallel default(none) shared(sum,a,m,wch) reduction(+:sum)
	Gomp_set_num_routines(4)	
	runtime.GOMAXPROCS(12)
	g:=Gomp_get_num_routines()
	for i := 0; i < g; i++ {
		go func(b int) {
			sum:=0
			fmt.Println("el valor de sum en el hilo",b,"es: ",sum)
			sum++
			wch <- sum
		}(i)
	}

	for i := 0; i < g; i++ {
		sum += <-wch
	}
	fmt.Println("el resultado de sum es: ",sum)
}
