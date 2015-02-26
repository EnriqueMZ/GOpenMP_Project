package main

import (
	"fmt"
	. "gomp_lib"
	//"io"
)

import "runtime"

// Comentario de prueba (  private(A,e) , firstprivate(B) reduction(+:f,g) , reduction(-: h, i) ). Porque estoy toqueteando comentarios
var n int
// var read io.Reader
// var i = false
// var fin bool = true
// var a, b, c, d int = 5, 2, 4, 7
// var test, pi, hello = false, 3.14, "Hello"
/* var e interface {
	// io.Reader
	 }*/
// var f [10]io.Reader
// var g *io.Reader
// var h map[io.Reader]io.Reader
// var p struct{ x, y float64 }
// var fun func(ident interface{}) string
// var ch chan int
// var ch_map chan map[io.Reader]io.Reader
/* var (
	i1			int
	j1				= true
	fin1			bool	= false
	a1, b1, c1		int	= 6, 8, 10
	test1, pi1, hello1		= true, 3.1416, "Bye"
	ch_map1			chan map[io.Reader]io.Reader
)*/
// var ()
var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func print_aux(ident int) int {
	var aux int = ident
	var a = func() int {
		var id int = aux
		return id
	}
	return a()
}
func main() {
	_init_numCPUs()
	fmt.Println("Inicio de la region paralela")
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < n+5; _i++ {
		go func(_routine_num int) {
			/*var (
				e	interface {
					io.Reader
				}
				f	[10]io.Reader
				g	*io.Reader
				h	map[io.Reader]io.Reader
				p	struct{ x, y float64 }
				fun	func(ident interface{}) string
				ch	chan int
				ch_map	map[io.Reader]io.Reader
				a	int
				b	int
				c	int
			)*/
			fmt.Println("Gouroutine:", _routine_num)

			var cont int = _routine_num
			for i := 0; i < 3; i++ {
				cont++
			}
			fmt.Println("cont =", cont)
			_barrier_0_bool <- true
		}(_i)
	}
	for _i := 0; _i < n+5; _i++ {
		<-_barrier_0_bool
	}

	fmt.Println("Fin de la region paralela")
	fmt.Println("Inicio de la region paralela")
	Gomp_set_num_routines(6)
	//var _barrier_0_bool = make(chan bool)  // WARNING!!! Revisar la creación de barreras.
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			fmt.Println("Gouroutine:", _routine_num)
			var cont int = _routine_num
			for i := 0; i < 3; i++ {
				cont++
			}
			fmt.Println("cont =", cont)
			_barrier_0_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_0_bool
	}

	fmt.Println("Fin de la region paralela")
}
