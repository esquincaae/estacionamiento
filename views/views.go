package views

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"carro/models"
	"image"
	_"image/png"
	"os"
)

var (
	fondo *pixel.Sprite
	imagenFondo  pixel.Picture
)

func cargarFondo() {
	archivo, err := os.Open("Assets/fondo.png")
	if err != nil {
		panic(err)
	}
	defer archivo.Close()

	img, _, err := image.Decode(archivo)
	if err != nil {
		panic(err)
	}

	imagenFondo = pixel.PictureDataFromImage(img)
	fondo = pixel.NewSprite(imagenFondo, imagenFondo.Bounds())
}

func DibujarEstacionamiento(ventana *pixelgl.Window, vehiculos []models.Vehiculo) {
	if fondo == nil {
		cargarFondo()
	}

	fondo.Draw(ventana, pixel.IM.Moved(ventana.Bounds().Center()))
	imd := imdraw.New(nil)
	imd.Color = colornames.White

	imd.Push(pixel.V(100, 500), pixel.V(700, 500))
	imd.Line(2)
	imd.Push(pixel.V(100, 100), pixel.V(700, 100))
	imd.Line(2)

	imd.Push(pixel.V(700, 100), pixel.V(700, 500))
	imd.Line(2)

	anchoEstacionamiento := 600.0 
	anchoCarril := anchoEstacionamiento / 10 

	for i := 0.0; i < 10.0; i++ {
		desplazamientoX := 100.0 + i*anchoCarril 
		imd.Push(pixel.V(desplazamientoX, 500), pixel.V(desplazamientoX, 350))
		imd.Line(2)  

		imd.Push(pixel.V(desplazamientoX, 250), pixel.V(desplazamientoX, 100))
		imd.Line(2)  
	}

	anchoVehiculo := anchoCarril / 4  
	altoVehiculo := anchoCarril / 4 

	for _, vehiculo := range vehiculos {
		imd.Color = colornames.Red

		esquinaSuperiorIzquierda := vehiculo.Posicion.Add(pixel.V(-anchoVehiculo, -altoVehiculo))
		esquinaSuperiorDerecha := vehiculo.Posicion.Add(pixel.V(anchoVehiculo, -altoVehiculo))
		esquinaInferiorIzquierda := vehiculo.Posicion.Add(pixel.V(-anchoVehiculo, altoVehiculo))
		esquinaInferiorDerecha := vehiculo.Posicion.Add(pixel.V(anchoVehiculo, altoVehiculo))

		imd.Push(esquinaSuperiorIzquierda, esquinaSuperiorDerecha, esquinaInferiorDerecha, esquinaInferiorIzquierda)
		imd.Polygon(0)
	}

	imd.Draw(ventana)
}
