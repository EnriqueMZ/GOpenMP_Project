package SAXPY_Bench

func Saxpy_parallel() {
	_init_numCPUs()
	var n int = 300000000
	var a float64 = 2
	x, y := Saxpy_init(n)
	
	var _barrier_1_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num + 0; i < (n+0)/1; i += _numCPUs {
				y[i] = a*x[i] + y[i]
			}
			_barrier_1_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_1_bool
	}

}