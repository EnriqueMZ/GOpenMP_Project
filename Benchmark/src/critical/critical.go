package main
//Ctrl+Shift+L miralo
import (
	"fmt"
	"runtime"
//	"sync"
)
const N=10
var a =5
var ch= make(chan int,1)  //canal para garantizar la exclusion mutua de la variable, mutex seria mas eficiente
//var mutex=sync.Mutex{}
var wch=make(chan int)
func main(){
	runtime.GOMAXPROCS(4)
	ch<-a
	for i:=0;i<N;i++{
		go func (){
			a1:=<-ch
			//mutex.Lock()
			a1++
			//mutex.Unlock()
			ch<-a1	
			wch<-0
		}()
		}
	for i:=0;i<N;i++{
		<-wch}
	a=<-ch
		fmt.Println("el resultado de a es: ",a)
	
	}
