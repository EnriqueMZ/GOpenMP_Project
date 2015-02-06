package main

import (
	"fmt"
	. "gomp_lib"
)

func main() {

	var cont int = 5
	var i int

	fmt.Println("Inicio de la region paralela")

	//pragma gomp parallel for

	for i := 0; i < 5; i++ {
		cont++
		fmt.Println("Gouroutine:", Gomp_get_routine_num(), " cont =", cont)
	}

	fmt.Println("Fin de la region paralela")
	fmt.Println("Valores de las variables fuera del bloque parallel:", i, cont)

}
