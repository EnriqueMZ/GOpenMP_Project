/*
 ==========================================================================
 Name        : gomp_lib.go
 Author      : Enrique Madridejos Zamorano
 Version     :
 Copyright   : Apache Licence Version 2.0
 Description : Biblioteca auxiliar para GOpenMP
 ==========================================================================
 */

package gomp_lib

import(
	"runtime"
	"time"
	)

var GOMP_NUM_ROUTINES int = runtime.NumCPU()

func Gomp_get_num_procs() int {
	return runtime.NumCPU()
	}

func Gomp_set_num_routines(N int){
	GOMP_NUM_ROUTINES = N
	}

func Gomp_get_num_routines() int {
	return GOMP_NUM_ROUTINES
	}

func Gomp_get_routine_num() int {
	return 0
	}
func Gomp_get_wtime() float64 {
	return float64(time.Now().UnixNano())/1e9
	}