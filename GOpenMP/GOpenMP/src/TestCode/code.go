package main

import (
	"fmt"
	. "gomp_lib"
	"io"
)

// Comentario de prueba (  private(A,e) , firstprivate(B) reduction(+:f,g) , reduction(-: h, i) ). Porque

var n int
var read io.Reader 
var i = false
var fin bool = true
var a, b, c, d int = 5, 2, 4, 7
var test, pi, hello = false, 3.14, "Hello"

var e interface{io.Reader}
var f [10]io.Reader
var g *io.Reader
var h map[io.Reader] io.Reader
var p struct { x, y float64 }
var fun func(ident interface{})string
var ch chan int
var ch_map chan map[io.Reader] io.Reader

var (
	i1                 int;
	j1                      = true;
	fin1               bool = false;
	a1, b1, c1         int  = 6, 8, 10;
	test1, pi1, hello1      = true, 3.1416, "Bye";
	ch_map1 chan map[io.Reader] io.Reader
)

var ()

func print_aux (ident int) int {
	var aux int = ident	
	var a = func() int {
		var id int = aux
		return id
		}
	return a()
	} 

func main() {
	
	fmt.Println("Inicio de la region paralela")

	//pragma    gomp   parallel  if (n>10) , num_threads(n + 5) private(a, b,c )   default(none), shared(d) private(e, f, g, h,p,fun , ch, ch_map)
	{
		fmt.Println("Gouroutine:", Gomp_get_routine_num())
		Gomp_set_num_routines(4)
		var cont int = Gomp_get_routine_num()
		for i := 0; i < 3; i++ {
			cont++
		}
		fmt.Println("cont =", cont)
	}

	fmt.Println("Fin de la region paralela")
	fmt.Println("Inicio de la region paralela")
	
	Gomp_set_num_routines(6)

	//pragma gomp parallel
	{
		fmt.Println("Gouroutine:", Gomp_get_routine_num())
		var cont int = Gomp_get_routine_num()
		for i := 0; i < 3; i++ {
			cont++
		}
		fmt.Println("cont =", cont)
	}

	fmt.Println("Fin de la region paralela")
}
