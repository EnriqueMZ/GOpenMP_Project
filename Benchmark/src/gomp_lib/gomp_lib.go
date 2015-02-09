package gomp_lib
import(
	"runtime"
	"os"
	. "strconv"
)
var GOMP_NUM_ROUTINES int = runtime.NumCPU()
func Gomp_set_num_routines(N int){
	os.Setenv("num_routines",Itoa(N))
}
func Gomp_get_num_routines() int {
	res,_:=Atoi(os.Getenv("num_routines"))
	return res
}

func Gomp_get_routine_num() int {
	return 0
}