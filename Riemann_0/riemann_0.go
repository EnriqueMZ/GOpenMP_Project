// Alberto Castaño

package main

import (
	"os"
	"strconv"
)

var step float64
var num_steps int

func main() {
	num_steps,_= strconv.Atoi(os.Args[1])
	var x, sum float64
	step = 1.0 / float64(num_steps)

	//pragma gomp parallel for reduction(+:sum)
	for i := 0; i < num_steps; i++ {
		x = (float64(i) + 0.5) * step
		sum = sum + 4.0/(1.0+x*x)

	}

	sum = step * sum
}