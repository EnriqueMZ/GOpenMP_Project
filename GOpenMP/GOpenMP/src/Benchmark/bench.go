/* package main

import (
	"fmt"
	. "gomp_lib"
)

func main() {

	fmt.Println("Inicio de la region paralela")	
	Gomp_set_num_routines(2)
	var cont1 int = 5
	
	//pragma gomp parallel num_threads(4)
	{
		var cont2 int = 5
		cont1++
		cont2++
		fmt.Println("Gouroutine:", Gomp_get_routine_num(), "cont1 =", cont1, "cont2 =", cont2)
	}
	
	fmt.Println("Fin de la region paralela")
}*/

package main

import (
	"fmt"
	. "gomp_lib"
)

import . "pragma_processor"
var pragma string = "//pragma gomp parallel num_threads(4)"
var prg_pro Pragma = ProcessPragma(pragma)

func main() {
	
	fmt.Println("Inicio de la region paralela")
	Gomp_set_num_routines(4)
	var cont1 int = 5
	
	var ch = make(chan int,1)
	ch <- cont1
	
	var done = make(chan bool)
	N := Gomp_get_num_routines()
	
	for i := 0; i < N; i++ {
		go func(id int) {
			var cont2 int = 5
			cont1 := <- ch
			cont1++
			ch <- cont1
			cont2++
			
			fmt.Println("Gouroutine:", id, "  Contador 1 =", cont1, "  Contador 2  =", cont2)
			
			done <- true
			
		}(i)
	}
	for i := 0; i < N; i++ {
		<-done
	}
	close(done)
	cont1 = <- ch
	close(ch)
	fmt.Println("Valor final de Contador 1 =", cont1)
	fmt.Println("Fin de la region paralela")
}
