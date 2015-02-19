package main 

import (
	"fmt"
)

func main() {
	
	var n int = 10;
	var a float64 = 2;
	x := []float64{1,2,3,4,5,6,7,8,9,10}
	y := []float64{1,2,3,4,5,6,7,8,9,10}
	
	fmt.Println("Vector x antes del parallel:", x)
	fmt.Println("Vector y antes del parallel:", y)
	
	//pragma gomp parallel for
	
	for i:= 0; i < n; i++ {
		y[i] = a * x[i] + y[i]
		}
	
	fmt.Println("Vector x despues del parallel:", x)
	fmt.Println("Vector y despues del parallel:", y)

}

