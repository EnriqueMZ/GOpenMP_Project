package gomp_lib

import(
	"runtime"
	)

var GOMP_NUM_ROUTINES int = runtime.NumCPU()

func Gomp_set_num_routines(N int){
	GOMP_NUM_ROUTINES = N
	}

func Gomp_get_num_routines() int {
	return GOMP_NUM_ROUTINES
	}

func Gomp_get_routine_num() int {
	return 0
	}

func Prueba_2(){ // Prueba para cambios
	} 