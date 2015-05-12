package Dot_Init

func Dot_Init_A() (int, [10000000]float64, [10000000]float64) {
	var a [10000000]float64
	var	b [10000000]float64
	var	n int = 10000000  // Size
	
	for i := 0; i < n; i++ {
		a[i] = float64(i) * 0.5
		b[i] = float64(i) * 2.0
		}
	
	return n, a, b
}

func Dot_Init_B(size int) ([]float64, []float64) {
	
	a := make([]float64, size)
	b := make([]float64, size)

	for i := 0; i < size; i++ {
		a[i] = float64(i) * 0.5
		b[i] = float64(i) * 2.0
	}
	return a, b
}
