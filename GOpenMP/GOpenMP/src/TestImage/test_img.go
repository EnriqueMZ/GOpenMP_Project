package main

import(
	"os"
	"flag"
	"log"
	"image"
	"image/jpeg"
	_"image/png"
	_"image/gif"
	"image/color"
	)

type img struct{
	h,w int
	c [][]color.RGBA
}

func (m img) At(x, y int) color.Color { return m.c[x][y] }
func (m img) ColorModel() color.Model { return color.RGBAModel }
func (m img) Bounds() image.Rectangle { return image.Rect(0, 0, m.h, m.w) }

var(
	
	img1=flag.String("img1","/home/acastano/prugo/src/imagen1.jpg","nombre de la foto")
	img2=flag.String("img2","/home/acastano/prugo/src/imagen3.jpg","nombre del archivo de salida")
)

func Create(imagen image.Image) img{
	a:=imagen.Bounds().Max.X-imagen.Bounds().Min.X
	b:=imagen.Bounds().Max.Y - imagen.Bounds().Min.Y
	c := make([][]color.RGBA,a)
	for i := range c {
		c[i] = make([]color.RGBA,b)
	}
	m := img{a,b,c}
	var tam int = m.h

	//pragma gomp parallel for
	for i:=0;i<tam;i++{
		for j:=0;j<m.w;j++{
			_a,_b,_c,_d:=imagen.At(i,j).RGBA()
			m.c[i][j].R, m.c[i][j].G, m.c[i][j].B, m.c[i][j].A = uint8(_a),uint8(_b),uint8(_c),uint8(_d)
		}
	}
	
	return m
}

func main(){
	//parsea los flags de entrada
	flag.Parse()

	//Crea el archivo de salida
	f,_:=os.Create(*img2)
	
	//abre el fichero a copiar
	f1,err:=os.Open(*img1)
	if err != nil {
	log.Fatal(err)
}
	//decodifica la imagen a copiar para obtener sus datos (tamaño y color)
	img,_,err:=image.Decode(f1)
	if err!=nil{
		log.Fatal(err)
	}
	
	img1:=Create(img)
	err = jpeg.Encode(f,img1,nil)
	err = jpeg.Encode(f, img,&jpeg.Options{jpeg.DefaultQuality})
}