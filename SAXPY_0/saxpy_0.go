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
	var a float32 = 2
	x:= make([]float32,n)
	y:= make([]float32,n)
	//pragma gomp parallel for
	for i:=0;i<n;i++{
		x[i]=float32(rand.Int())
		y[i]=float32(rand.Int())
	}
	
	//pragma gomp parallel for
	
	for i:= 0; i < n; i++ {
		y[i] = a * x[i] + y[i]
		}

}

