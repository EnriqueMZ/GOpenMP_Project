// Alberto Castaño

package main

import(
	"os"
	"strings"
	. "io/ioutil"
	"fmt"
	"gomp_lib"
	"strconv"
	)

func duplicate(a string,n int) string{
	for i:=0;i<n;i++{
		a+=" "+a
	}
	return a
}

func ReadAndSplit(s string,n int) ([]string,float64){
	a:=gomp_lib.Gomp_get_wtime()
	b,_:=ReadFile(s)
	b_0:=duplicate(string(b),n)
	f:=strings.Split(b_0," ")
	a=gomp_lib.Gomp_get_wtime()-a
	return f,a
}

func AnalizarTexto(s []string) (map[string]int,float64){
	a:=gomp_lib.Gomp_get_wtime()
	elements := make(map[string]int)
	var m int=len(s)
	//pragma gomp parallel for 
	for i:=0;i<m;i++{
			elements[s[i]]++
	}
	a=gomp_lib.Gomp_get_wtime()-a
	return elements,a
}


func main(){
	a:=os.Args[1]
	b,_:=strconv.Atoi(os.Args[2])
	s,t1:=ReadAndSplit(a,b)
	m,t2:=AnalizarTexto(s)
	fmt.Println(m["Harry"])
	fmt.Println("tiempo en leer y duplicar: ",t1)
	fmt.Println("tiempo en analizar texto: ",t2)
}