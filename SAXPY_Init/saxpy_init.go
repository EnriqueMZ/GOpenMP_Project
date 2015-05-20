/*
==================================================================
 SAXPY: Init para inicializar los vectores
==================================================================
*/

package SAXPY_Init

import "math/rand"

func Saxpy_init(size int) ([]float64, []float64) {

	x := make([]float64, size)
	y := make([]float64, size)

	for i := 0; i < size; i++ {
		x[i] = float64(rand.Int())
		y[i] = float64(rand.Int())
	}
	return x, y
}
