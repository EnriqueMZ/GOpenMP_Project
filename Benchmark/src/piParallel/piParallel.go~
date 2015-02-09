package main

import (
	"fmt"
	. "gomp_lib"
	"runtime"
)

var step, x float64
var num_steps = 10000000000

func main() {
	runtime.GOMAXPROCS(4)
	var pi float64
	step = 1.0 / float64(num_steps)
	//pragma gomp parallel num_threads(4)
	ch := make(chan float64)
	Gomp_set_num_routines(4)
	g := Gomp_get_num_routines()
	for i := 0; i < g; i++ {
		go func(tid int) {
			sum := 0.0
			x:=0.0
			//pragma gomp for
			for j := tid; j < num_steps; j = j + g {
				x = (float64(j) + 0.5) * step
				sum = sum + 4.0/(1.0+x*x)
			}
			ch <- sum * step
		}(i)

	}
	for i := 0; i < Gomp_get_num_routines(); i++ {
		pi += <-ch
	}
	fmt.Println("La aproximacion de numero Pi con ", num_steps, "iteraciones es: ", pi)
}
