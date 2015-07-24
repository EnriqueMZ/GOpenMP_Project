package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	output = flag.String("out", "mandelbrot_1.jpeg", "name of the output image file")
	height = flag.Int("h", 4096, "height of the output image in pixels")
	width  = flag.Int("w", 4096, "width of the output image in pixels")
)

type img struct {
	h, w int
	m    [][]color.RGBA
}

var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}

var init_p time.Time

func (m *img) At(x, y int) color.Color { return m.m[x][y] }
func (m *img) ColorModel() color.Model { return color.RGBAModel }
func (m *img) Bounds() image.Rectangle { return image.Rect(0, 0, m.h, m.w) }
func Create(h, w int) image.Image {
	c := make([][]color.RGBA, h)
	for i := range c {
		c[i] = make([]color.RGBA, w)
	}
	m := img{h, w, c}
	var mm int = len(m.m)
	init_p = time.Now()
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num + 0; i < (mm+0)/1; i += _numCPUs {
				for j := 0; j < mm; j++ {
					fillPixel(&m, i, j)
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
func fillPixel(m *img, i, j int) { // normalized from -2.5 to 1
	xi := 3.5*float64(i)/float64(m.w) - 2.5
	// normalized from -1 to 1
	yi := 2*float64(j)/float64(m.h) - 1
	const maxI = 1000
	x, y := 0., 0.
	for i := 0; (x*x+y*y < 4) && i < maxI; i++ {
		x, y = x*x-y*y+xi, 2*x*y+yi
	}
	paint(&m.m[i][j], x, y)
}
func paint(c *color.RGBA, x, y float64) {
	n := byte(x * y)
	c.R, c.G, c.B, c.A = n, n, n, 255
}
func main() {
	init := time.Now()
	_init_numCPUs()
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	// open a new file
	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	// create the image
	img := Create(*height, *width)
	// and encoding it
	fimt := filepath.Ext(*output)
	switch fimt {
	case ".png":
		err = png.Encode(f, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, img, nil)
	case ".gif":
		err = gif.Encode(f, img, nil)
	default:
		err = errors.New("unkwnown format " + fimt)
	}
	// unless you can't
	if err != nil {
		log.Fatal(err)
	}
	fin_p := time.Since(init_p).Seconds()
	fin := time.Since(init).Seconds()
	fmt.Println(_numCPUs, ",", fin_p, ",", fin)
}
