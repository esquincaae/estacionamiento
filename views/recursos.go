package views

import (
	"image"
	"os"
	"github.com/faiface/pixel"
	_"image/png"
)

var (
	fondo *pixel.Sprite
	imagenFondo  pixel.Picture
	imagenCarro *pixel.Sprite
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

func cargarCarro() {
    archivo, err := os.Open("Assets/carro.png")
    if err != nil {
        panic(err)
    }
    defer archivo.Close()

    img, _, err := image.Decode(archivo)
    if err != nil {
        panic(err)
    }

    pic := pixel.PictureDataFromImage(img)
    imagenCarro = pixel.NewSprite(pic, pic.Bounds())
}
