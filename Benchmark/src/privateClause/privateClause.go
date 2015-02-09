package main

import (
	"fmt"
	. "gomp_lib"
)


var a, b, c = 1, 2, 3

func main() {
	fmt.Println("valores de a, b y c antes de la region paralela: ", a, b, c)
	//pragma gomp parallel(5) private (a,b,c)
	ch := make(chan int)
	Gomp_set_num_routines(5)
	for i := 0; i < Gomp_get_num_routines(); i++ {
		go func(y int) { //lanzamos tantas rutinas como numero de threads nos pida la directiva
			var a, b, c int
			fmt.Println("valores de a, b y c en region paralela con private :", a, b, c, "Soy la rutina: ", Gomp_get_routine_num())
			//fin de la region paralela, devolviendo valores iniciales a las variables
			ch <- 0 //indicamos al canal la terminacion de la rutina
		}(i)
	}
	for i := 0; i < Gomp_get_num_routines(); i++ {
		<-ch
	}
}
