package SAXPY_Bench

import "math/rand"
import "runtime"

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
	}

func Saxpy_init(size int) ([]float32, []float32) {

	x := make([]float32, size)
	y := make([]float32, size)

	for i := 0; i < size; i++ {
		x[i] = float32(rand.Int())
		y[i] = float32(rand.Int())
	}
	return x, y
}
