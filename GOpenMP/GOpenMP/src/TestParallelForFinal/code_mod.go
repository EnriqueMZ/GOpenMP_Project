package main

import (
	"fmt"
	//. "gomp_lib"
	"runtime" // Incluido
)

func main() {

	_numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(_numCPUs)

	var cont int = 5
	var i int

	fmt.Println("NumCPUs: ", _numCPUs)
	fmt.Println("Inicio de la region paralela")

	/*
		//pragma gomp parallel for


			for i := 0; i < 5; i++ {
				cont++
				fmt.Println("Gouroutine:", Gomp_get_routine_num(), " cont =", cont)
			}

		fmt.Println("Fin de la region paralela")
		fmt.Println("Valores de las variables fuera del bloque parallel:", i, cont)

		fmt.Println("Inicio de la region paralela")
	*/

	//pragma gomp parallel for

	var _barrier = make(chan bool)
	for _i := 0; _i < 8; _i++ { // Hilos.
		go func(_routine_num int) {
			for _j := _routine_num; _j < 10; _j += 8 { // Iteraciones + Hilos.
				cont++
				fmt.Println("Gouroutine:", _routine_num, " cont =", cont)
			}
			_barrier <- true
		}(_i)
	}
	for _i := 0; _i < 8; _i++ { // Hilos
		<-_barrier
	}

	fmt.Println("Fin de la region paralela")
	fmt.Println("Valores de las variables fuera del bloque parallel:", i, cont)

}
