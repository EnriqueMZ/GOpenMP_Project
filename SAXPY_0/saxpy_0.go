// Alberto Castaño

package main 

import (
	"os"
	"strconv"
	"math/rand"
)
func main() {
	var n int
	n,_= strconv.Atoi(os.Args[1])
	var a float64 = 2
	x:= make([]float64,n)
	y:= make([]float64,n)
	//pragma gomp parallel for
	for i:=0;i<n;i++{
		x[i]=float64(rand.Int())
		y[i]=float64(rand.Int())
	}
	
	//pragma gomp parallel for
	
	for i:= 0; i < n; i++ {
		y[i] = a * x[i] + y[i]
		}

}

