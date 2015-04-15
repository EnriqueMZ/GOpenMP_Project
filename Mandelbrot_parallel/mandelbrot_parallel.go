package main

import (
	"errors"
	"flag"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"fmt"
)

var (
	output = flag.String("out", "mandelbrot_1.jpeg", "name of the output image file")
	height = flag.Int("h", 1024, "height of the output image in pixels")
	width  = flag.Int("w", 1024, "width of the output image in pixels")
)

type img struct {
	h, w int
	m    [][]color.RGBA
}

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func (m *img) At(x, y int) color.Color { return m.m[x][y] }
func (m *img) ColorModel() color.Model { return color.RGBAModel }
func (m *img) Bounds() image.Rectangle { return image.Rect(0, 0, m.h, m.w) }

func Create(h, w int) image.Image {
	
	c := make([][]color.RGBA, h)
	for i := range c {
		c[i] = make([]color.RGBA, w)
	}
	
	count := make([][]int, h)
	for i := range c {
		count[i] = make([]int, w)
	}

	m := img{h, w, c}
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num + 0; i < (w+0)/1; i += _numCPUs {
				for j := 0; j < h; j++ {
					var i int
					fillPixel(&m, i, j, count)
				}
			}
			_barrier_0_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_0_bool
	}

	return &m
}
func fillPixel(m *img, i int, j int, count [][]int) {
	var x_max float64 = 1.25
	var x_min float64 = -2.25
	var y_max float64 = 1.75
	var y_min float64 = -1.75

	x := (float64(j-1)*x_max + float64(m.w-j) + x_min) / float64(m.w-1)
	y := (float64(i-1)*y_max + float64(m.h-i) + y_min) / float64(m.h-1)
	fmt.Println("Valores de i, j", i ,j)
	
 	count[i][j] = 0;

	const maxI = 2000

	x1 := x
	y1 := y

	for k := 1; k <= maxI; k++ {
		x2 := x1*x1 - y1*y1 + x
		y2 := 2*x1*y1 + y

		if x2 < -2.0 || 2.0 < x2 || y2 < -2.0 || 2.0 < y2 {
			count[i][j] = k
			break
		}
		x1 = x2
		y1 = y2
	}
	if (count[i][j] % 2) == 1 {
		paint(&m.m[i][j], 255)
	} else {
		c := (255.0 * math.Sqrt(math.Sqrt(math.Sqrt((float64(count[i][j]) / float64(maxI))))))
		paint_complex(&m.m[i][j], c)
	}
}

func paint(c *color.RGBA, k float64) {
	n := byte(k)
	c.R, c.G, c.B, c.A = n, n, n, 255
}
func paint_complex(c *color.RGBA, k float64) {
	a := byte(k)
	b := byte(3 * k / 5)
	c.R, c.G, c.B, c.A = b, b, a, 255
}

func main() {
	_init_numCPUs()
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	// open a new file estoy toqueteando comentarios
	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	// create the image estoy toqueteando comentarios
	img := Create(*height, *width)
	// and encoding it estoy toqueteando comentarios
	fmt := filepath.Ext(*output)
	switch fmt {
	case ".png":
		err = png.Encode(f, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, img, nil)
	case ".gif":
		err = gif.Encode(f, img, nil)
	default:
		err = errors.New("unkwnown format " + fmt)
	}
	// unless you can't estoy toqueteando comentarios
	if err != nil {
		log.Fatal(err)
	}
}
