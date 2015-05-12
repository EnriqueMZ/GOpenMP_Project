package Dot_Bench

//import "fmt" 

func Dot_serial_A() {
	var sum float64 = 0
	
	n, a, b := Dot_Init_A() 
	
	for i := 0; i < n; i++ {
		sum += a[i] * b[i]
	}
	//fmt.Println("Serial A result: ", sum)
}

func Dot_serial_B() {
	var sum float64 = 0
	var n int = 300000000
	
	a, b := Dot_Init_B(n) 
	
	for i := 0; i < n; i++ {
		sum += a[i] * b[i]
	}
	//fmt.Println("Serial B result: ", sum)
}