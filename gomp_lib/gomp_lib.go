/*
 ==========================================================================
 Name        : gomp_lib.go
 Author      : Enrique Madridejos Zamorano
 Version     :
 Copyright   : Licensed under the Apache License, Version 2.0 (the "License");
   			   you may not use this file except in compliance with the License.
               You may obtain a copy of the License at

               http://www.apache.org/licenses/LICENSE-2.0

               Unless required by applicable law or agreed to in writing, software
               distributed under the License is distributed on an "AS IS" BASIS,
               WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
               See the License for the specific language governing permissions and
               limitations under the License.
               
 Description : Auxiliar GOpenMP library
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