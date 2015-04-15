// Alberto Castaño

package main


import (
	"fmt"
	"strconv"
	"strings"
	"os"
	. "io/ioutil"
	)


	func main(){

		np,_:=strconv.Atoi(os.Args[1])
		nt,_:=strconv.Atoi(os.Args[2])
		c,_:=ReadFile("resultados.txt")
		f:=strings.Split(string(c),"\n")

		resultado:= make([]float64,np)
		var ch=make(chan int)
		for i:=0;i<np;i++{
			go func(a int) {
				res:=0.0
			for i := a; i <(a+1); i ++ {
				d,_:=strconv.ParseFloat(f[i],64)
				res+=d
			}
			resultado[a]=res/float64(nt)
			ch<-0
		}(i)
		}
		for i:=0;i<19;i++{
		<-ch
		}
		for i:=0;i<19;i++{
		fmt.Println(resultado[i])
	}

	}