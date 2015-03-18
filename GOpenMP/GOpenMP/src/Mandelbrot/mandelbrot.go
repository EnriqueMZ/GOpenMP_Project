package main

import (	
	"flag"
	"math"
	"os"
	"runtime"
	"fmt"
)



var (
	output = flag.String("out", "mandelbrot_1.ppm", "name of the output image file")
	height = flag.Int("h", 700, "height of the output image in pixels")
	width  = flag.Int("w", 700, "width of the output image in pixels")
)

func Create(h, w int) ([][]int, [][]int, [][]int) {
	
	count := make([][]int, h)
	for i := range count {
		count[i] = make([]int, w)
	}
	
	r := make([][]int, h)
	for i := range r {
		r[i] = make([]int, w)
	}
	
	g := make([][]int, h)
	for i := range g {
		g[i] = make([]int, w)
	}
	
	b := make([][]int, h)
	for i := range b {
		b[i] = make([]int, w)
	}

	var i int
	var j int
	var k int
	var x float64
	var x1 float64
	var x2 float64
	var y float64
	var y1 float64
	var y2 float64
	const maxI int = 2000
	var x_max float64 = 1.25
	var x_min float64 = -2.25
	var y_max float64 = 1.75
	var y_min float64 = -1.75
	var res int

	//pragma gomp parallel for shared (count, maxI, x_max, x_min, y_max, y_min ) private ( i, j, k, x, x1, x2, y, y1, y2 )
	for i = 0; i < w; i++ {
		for j = 0; j < h; j++ {

			x = (float64(j-1)*x_max + float64(w-j) + x_min) / float64(w-1)
			y = (float64(i-1)*y_max + float64(h-i) + y_min) / float64(h-1)

			count[i][j] = 0

			x1 = x
			y1 = y

			for k = 1; k <= maxI; k++ {
				x2 = x1 * x1 - y1 * y1 + x
				y2 = 2 * x1 * y1 + y

				if x2 < -2.0 || 2.0 < x2 || y2 < -2.0 || 2.0 < y2 {
					count[i][j] = k
				}
				x1 = x2
				y1 = y2
			}
			if (count[i][j] % 2) == 1 {
				r[i][j] = 255
        		g[i][j] = 255
        		b[i][j] = 255				
			} else {
				res = int(255.0 * math.Sqrt(math.Sqrt(math.Sqrt( float64(count[i][j]) / float64(maxI) ))))
				r[i][j] = 3 * res / 5;
        		g[i][j] = 3 * res / 5;
        		b[i][j] = res;
			}

		}
	}
	return r, g, b
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func i4_min( i1 int, i2 int ) int {

  var value int

  if  i1 < i2 {
    value = i1;
  }else{
    value = i2;
  }
  return value;
}

func main() {
	
	var jhi int
	
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	//create the image
	r, g ,b := Create(*height, *width)

	// open a new file
	f, err := os.Create(*output)
	check(err)
	
	defer f.Close()
	
	fmt.Fprint(f, "P3\n")
	fmt.Fprintf(f, "%d  %d\n", *height, *width )
	fmt.Fprintf(f, "%d\n", 255 )
	
	for  i := 0; i < *height; i++ {
    	for jlo := 0; jlo < *width; jlo = jlo + 4 {
      	jhi = i4_min( jlo + 4, *width )
      	for  j := jlo; j < jhi; j++ {
        	fmt.Fprintf(f, "  %d  %d  %d", r[i][j], g[i][j], b[i][j])
      	}
      	fmt.Fprint(f, "\n")
    	}
    	
  	}
}
