/*
==================================================================
 SAXPY: Init para inicializar los vectores
==================================================================
*/

package SAXPY_Init

import "math/rand"

func Saxpy_init(size int) ([]float32, []float32) {

	x := make([]float32, size)
	y := make([]float32, size)

	for i := 0; i < size; i++ {
		x[i] = float32(rand.Int())
		y[i] = float32(rand.Int())
	}
	return x, y
}
