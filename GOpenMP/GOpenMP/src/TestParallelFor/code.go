package main

import (
	"fmt"
	. "gomp_lib"
)

func main() {

	var sum1 int = 0
	var sum2 int = 0
	var prod float64 = 0
	var res float64 = 1000
	var cont int = 0

	fmt.Println("Inicio de la region paralela")

	//pragma gomp parallel for reduction(+:sum1, sum2) reduction(*:prod) private(cont1)

	for i := 0; i < 10; i += 2 {
		sum1 += 1
		sum2 += 2
		prod *= 2
		res -= 2
		cont++
		fmt.Println("Gouroutine:", Gomp_get_routine_num(), " cont =", cont)
	}

	fmt.Println("Fin de la region paralela")
	fmt.Println("Valores de la variable fuera del bloque parallel:", sum1, sum2, prod)

}
