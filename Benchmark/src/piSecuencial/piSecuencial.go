package main

import (
	"fmt"
	"runtime"
)

var step float64
var num_steps = 10000000000

func main() {
	runtime.GOMAXPROCS(4)
	var x, pi, sum float64
	step = 1.0 / float64(num_steps)
	//pragma gomp parallel num_threads(100)

	//pragma gomp for
	for i := 0; i < num_steps; i++ {
		x = (float64(i) + 0.5) * step
		sum = sum + 4.0/(1.0+x*x)

	}

	pi = step * sum
	fmt.Println("La aproximacion de numero Pi con ", num_steps, "iteraciones es: ", pi)
}