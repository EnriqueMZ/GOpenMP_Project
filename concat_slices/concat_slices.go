package main

import (
	"fmt"
	"reflect"
)

//var listMap map[int] []string

func incBucle() int {
	var cont int = 0
Loop:
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			cont++
			break Loop
		}
	}
	return cont
}

func notFound(a string, varList []string) (bool, string) {
	var res bool = true
	var v string = ""
	if len(varList) == 0 {
		res = false
		v = a
	} else {
		for i := range varList {
			if a == varList[i] {
				res = false
				break
			}
			if i == len(varList)-1 {
				v = a
			}
		}
	}
	return res, v
}

func notFoundList(pragma, varList []string) (bool, string) {
	var res bool
	var v string
	if len(pragma) == 0 {
		res = true
		v = ""
	} else {
		for i := range pragma {
			res, v = notFound(pragma[i], varList)
			if res {
				break
			}
		}
	}
	return res, v
}

func repeatIn(a []string) (bool, string) {
	var res bool = false
	var elem = ""
	for i := 0; i < len(a)-1; i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] == a[j] {
				res = true
				elem = a[i]
			}
		}
	}
	return res, elem
}

func repeat(a, b []string) bool {
	var res bool = false
	for i := range b {
		for j := range a {
			if a[j] == b[i] {
				res = true
			}
		}
	}
	return res
}

func repeatOne(a []string, b [][]string) bool {
	var res bool = false
	for i := range b {
		res = repeat(a, b[i])
	}
	return res
}

func concat(a, b []string) []string {
	for i := range b {
		a = append(a, b[i])
	}
	return a
}

func tiping (ident interface{}) reflect.Type {
	return reflect.TypeOf(ident)
	}

func inc (a int) int {
	a++
	return a
	}

func main() {
	//emp := []string{}
	//a :=
	//b := []string{"d", "e", "f"}
	//c := []string{"g", "h", "i"}
	//d := []string{"A", "B", "C"}
	//e := [][]string{b,c,d}
	 
	var ts1 = 1
	var ts2 = 1.2
	var ts3 = 1.2i
	var ts4 = "Hola"
	var ts5 = false
	var ts6, ts7, ts8 = 1, 1.3, 'a'
	var ts9 = 3.36e3
	
	var ts10 = func (ident interface{}) string {
	return reflect.TypeOf(ident).String()
	}
	
	var ts11 = ts10
	ts11 = nil
	
	var ts12 = 2 + 5
	
	var ts13 = reflect.Zero(reflect.TypeOf(ts12)).Interface() // Funciona
	
	//var ts13 = reflect.Zero(reflect.TypeOf(ts12)).Interface() // Funciona
	
	//ts14 := ts12 + ts13
	
	fmt.Println(ts1, ts2, ts3, ts4, ts5, ts9, ts10, ts11, ts12, ts13)
	fmt.Println(tiping(ts1))
	fmt.Println(tiping(ts2))
	fmt.Println(tiping(ts3))
	fmt.Println(tiping(ts4))
	fmt.Println(tiping(ts5))
	fmt.Println(tiping(ts6))
	fmt.Println(tiping(ts7))
	fmt.Println(tiping(ts8))
	fmt.Println(tiping(ts9))
	fmt.Println(tiping(ts10))
	fmt.Println(tiping(ts11))
	fmt.Println(tiping(ts12))
	fmt.Println(tiping(ts13))
	
	//variable := "d"
	//pragma := []string{"a", "b", "c"}
	//varList := []string{"a", "c"}
	//fmt.Println(notFound(variable, varList))
	//fmt.Println(notFoundList(pragma, varList))
	//fmt.Println(repeat(a, d))
	//fmt.Println(repeatOne(a, e))
	//fmt.Println(concat(a, b), b)
	//e := concat(b, d)
	//fmt.Println(e, repeat(a, e))
	//f := []string{"a", "b", "c", "d"}
	//fmt.Println(repeatIn(f))
	//fmt.Println(incBucle())
	/*listMap = make(map[int] []string)
	for key := 0; key <= 3; key++ {
		listMap[key] = emp
		}
	fmt.Println(listMap)
	listMap[0]= a
	listMap[1]= b
	listMap[2]= c
	listMap[3]= d
	fmt.Println(listMap)*/
}
