package SAXPY_Bench

func Saxpy_serial() {
	
	var n int = 300000000
	var a float32 = 2
	x, y := Saxpy_init(n)
	
	//pragma gomp parallel for
	
	for i:= 0; i < n; i++ {
		y[i] = a * x[i] + y[i]
		}
}