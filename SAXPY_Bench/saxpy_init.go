package SAXPY_Bench

import "math/rand"
import "runtime"

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
	}

func Saxpy_init(size int) ([]float64, []float64) {

	x := make([]float64, size)
	y := make([]float64, size)

	for i := 0; i < size; i++ {
		x[i] = float64(rand.Int())
		y[i] = float64(rand.Int())
	}
	return x, y
}
