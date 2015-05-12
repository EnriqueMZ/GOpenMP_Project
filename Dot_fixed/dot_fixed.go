package main 

import (
	"fmt"
	"time"
	"Dot_Init" 
	)

func Dot_serial_A() {
	var sum float64 = 0
	
	n, a, b := Dot_Init.Dot_Init_A() 
	init := time.Now()
	for i := 0; i < n; i++ {
		sum += a[i] * b[i]
	}
	fin := time.Since(init)
	fmt.Println("Time: ", fin)
	fmt.Println("Serial A result: ", sum)
}

func Dot_serial_B() {
	var sum float64 = 0
	var n int = 300000000
	init := time.Now()
	a, b := Dot_Init.Dot_Init_B(n) 
	
	for i := 0; i < n; i++ {
		sum += a[i] * b[i]
	}
	fin := time.Since(init)
	fmt.Println("Time: ", fin)
	fmt.Println("Serial B result: ", sum)
}

func main() {
	//Dot_serial_A() 
	Dot_serial_B()
}

