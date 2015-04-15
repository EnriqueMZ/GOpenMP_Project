// Alberto Castaño

package main

import (
	"os"
	"flag"
	"log"
	"image"
	"image/jpeg"
	//_"image/png"
	//_"image/gif"
	"image/color"
)

import "runtime"
//ey
//Estrucuta de imagen
type img struct {
	h, w	int
	c	[][]color.RGBA
}

//Metodos que se necesitan para manejar imagen
var _numCPUs = runtime.NumCPU()

func _init_numCPUs() {
	runtime.GOMAXPROCS(_numCPUs)
}
func (m img) At(x, y int) color.Color	{ return m.c[x][y] }
func (m img) ColorModel() color.Model	{ return color.RGBAModel }
func (m img) Bounds() image.Rectangle	{ return image.Rect(0, 0, m.h, m.w) }

/*
El maximo cogera un array de valores y de ellos devolvera el maximo
no como la mierda de Max(a float64,b float64)
*/
func Max(a []uint8) uint8 {
	var max uint8 = a[0]
	for i := 0; i < len(a)-1; i++ {
		if max < a[i+1] {
			max = a[i]
		}
	}
	return max
}

/*
Again
*/
func Min(a []uint8) uint8 {
	var min uint8 = a[0]
	for i := 0; i < len(a)-1; i++ {
		if min > a[i+1] {
			min = a[i]
		}
	}
	return min
}

/*
Un par de flags que permitan al programador elegir imagen de entrada y nombre
de nueva imagen de saldia.
*/
var (
	img1	= flag.String("img1", "../../joker.jpg", "nombre de la foto")
	img2	= flag.String("img2", "../../imagen3.jpg", "nombre del archivo de salida")
)

/*
Este metodo se encarga de crear una nueva imagen de tipo img y haremos los
filtros necesarios para el resultado(en este caso paso a escala de grises)
*/
func Create(imagen image.Image) img {	/*
		Los datos de la nueva imagen los seleccionamos en funcion de la
		de entrada.
	*/
	a := imagen.Bounds().Max.X - imagen.Bounds().Min.X
	b := imagen.Bounds().Max.Y - imagen.Bounds().Min.Y
	c := make([][]color.RGBA, a)
	for i := range c {
		c[i] = make([]color.RGBA, b)
	}
	/*
		Creamos la nueva imagen
	*/
	m := img{a, b, c}
	//cosas del preprocesador
	var mm int = m.h
	var _barrier_0_bool = make(chan bool)
	for _i := 0; _i < _numCPUs; _i++ {
		go func(_routine_num int) {
			var ()
			for i := _routine_num + 0; i < (mm+0)/1; i += _numCPUs {
				for j := 0; j < m.w; j++ {	/*
						Sacamos los colores de la imagen de entrada
					*/
					_r, _g, _b, _a := imagen.At(i, j).RGBA()
					/*
						cd es el valor devuelto por el algoritmo que elijamos para
						modificar el color de la imagen
					*/
					//cd:=(Max([]uint8{uint8(_r),uint8(_g),uint8(_b)})+Min([]uint8{uint8(_b),uint8(_g),uint8(_b)}))/2
					//cd:=(uint8(_a)+uint8(_b)+uint8(_c))/3
					cd := float64(uint8(_r))*0.30 + float64(uint8(_g))*0.59 + float64(uint8(_b))*0.11
					/*
						Seleccionamos color de cada pixel de la nueva imagen, valores uint8
						aunque los devueltos por la imagen original son de 32
					*/
					m.c[i][j].R, m.c[i][j].G, m.c[i][j].B, m.c[i][j].A = uint8(cd), uint8(cd), uint8(cd), uint8(_a)
				}
			}
			_barrier_0_bool <- true
		}(_i)
	}
	for _i := 0; _i < _numCPUs; _i++ {
		<-_barrier_0_bool
	}

	return m
}
func main() {
	_init_numCPUs()
	//parsea los flags de entrada
	flag.Parse()
	//Crea el archivo de salida
	f, _ := os.Create(*img2)
	//abre el fichero a copiar
	f1, err := os.Open(*img1)
	if err != nil {
		log.Fatal(err)
	}
	//decodifica la imagen a copiar para obtener sus datos (tamaño y color)
	img, _, err := image.Decode(f1)
	if err != nil {
		log.Fatal(err)
	}
	/*
		Llamada a la funcion que cumplira la funcion de modificar la imagen pixel
		a pixel e ir guardando dichos pixel en la matriz de colores de la nueva
		imagen.
	*/
	img1 := Create(img)
	err = jpeg.Encode(f, img1, nil)
	err = jpeg.Encode(f, img, &jpeg.Options{100})
}
