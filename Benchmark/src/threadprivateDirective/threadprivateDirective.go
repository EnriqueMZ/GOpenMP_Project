package main

import (
	"fmt"
	. "gomp_lib"
)

var a, b, c = 1, 2, 3
var tid int

func main() {

	//pragma gomp threadprivate(a,b,c)

	ba := [3]int{a, b, c}

	//pragma gomp parallel num_threads(5)

	Gomp_set_num_routines(5) //establecemos las 5 rutinas que se lanzaran
	vch := make([]chan interface{}, Gomp_get_num_routines())
	for i := 0; i < Gomp_get_num_routines(); i++ {

		vch[i] = make(chan interface{}, 3)
	}

	ch := make(chan int) //canal de espera

	for i := 0; i < Gomp_get_num_routines(); i++ {
		go func(q int) {
			tid = q
			a = ba[0]
			b = ba[1]
			c = ba[2]
			if tid == 1 {
				a += 10
			}
			a += 10
			b += 10
			c += 10
			vch[tid] <- a
			vch[tid] <- b
			vch[tid] <- c
			fmt.Println("Thread: ", tid, "hola", a, b, c)
			ch <- 0
		}(i)

	}

	for i := 0; i < Gomp_get_num_routines(); i++ {
		<-ch
	}

	//pragma gomp parallel num_threads(8) copyin(a)

	ch1 := make(chan int)    //canal de espera
	Gomp_set_num_routines(8) //establecemos las 5 rutinas que se lanzaran

	vch1 := make([]chan interface{}, Gomp_get_num_routines())
	copy(vch1, vch)
	for i := 5; i < Gomp_get_num_routines(); i++ {
		vch1[i] = make(chan interface{}, 3)
	}
	
	
	fmt.Println("**************************")
	fmt.Println("Master thread dando cañita")
	fmt.Println("**************************")
	
	
	
	for i := 0; i < Gomp_get_num_routines(); i++ {
		go func(q int) {
			tid = q
			if tid < 5 {
				a = (<-vch1[tid]).(int)
				a = ba[0]
				b = (<-vch1[tid]).(int)
				c = (<-vch1[tid]).(int)
			} else {
				a = ba[0] //la clausula copyin obliga a coger la copia inicial de threadprivate
				b = ba[1]
				c = ba[2]
			}

			fmt.Println("Thread: ", tid, "hola", a, b, c)
			vch1[tid] <- a
			vch1[tid] <- b
			vch1[tid] <- c
			ch1 <- 0
		}(i)

	}
	for i := 0; i < Gomp_get_num_routines(); i++ {
		<-ch1
	}
}
