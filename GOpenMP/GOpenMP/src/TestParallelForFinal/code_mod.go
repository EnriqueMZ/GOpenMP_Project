package main

import (
	"fmt"
	//. "gomp_lib"
)

import "runtime" // Incluido

func main() {

	_numCPUs := runtime.NumCPU() // Incluido
	runtime.GOMAXPROCS(_numCPUs) // Incluido

	var cont int = 5
	var i int

	fmt.Println("NumCPUs: ", _numCPUs)
	fmt.Println("Inicio de la region paralela")

	/*
		//pragma gomp parallel for


			for i := 0; i < 10; i+=2 {
				cont++
				fmt.Println("Gouroutine:", Gomp_get_routine_num(), " cont =", cont)
			}

		fmt.Println("Fin de la region paralela")
		fmt.Println("Valores de las variables fuera del bloque parallel:", i, cont)

		fmt.Println("Inicio de la region paralela")
	*/

	//pragma gomp parallel for

	var _barrier = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ { // Hilos.
		go func(_routine_num int) {
			for _i := _routine_num; _i < 5; _i += _numCPUs { // Iteraciones reales + Hilos.
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
