package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
)

const (
	anchoVentana  = 800.0
	altoVentana   = 600.0
	espacios      = 5
	tamanoEspacio = anchoVentana / float64(espacios+2)
	tamanoAuto    = tamanoEspacio * 0.6
)

type Estacionamiento struct {
	espacios   int
	mu         sync.Mutex
}

func NuevoEstacionamiento(capacidad int) *Estacionamiento {
	return &Estacionamiento{
		espacios: capacidad,
	}
}

func (e *Estacionamiento) Entrar() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.espacios > 0 {
		pos := espacios - e.espacios
		e.espacios--
		return pos
	}
	return -1
}

func (e *Estacionamiento) Salir(pos int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.espacios++
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Estacionamiento",
		Bounds: pixel.R(0, 0, anchoVentana, altoVentana),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	e := NuevoEstacionamiento(espacios)

	autos := make([]bool, espacios)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 500)
			pos := e.Entrar()
			if pos != -1 {
				fmt.Println("Vehículo", i, "entró y se estacionó en el espacio", pos)
				autos[pos] = true
				time.Sleep(time.Second * 2)
				e.Salir(pos)
				autos[pos] = false
				fmt.Println("Vehículo", i, "salió del espacio", pos)
			} else {
				fmt.Println("Vehículo", i, "no encontró espacio")
			}
		}
	}()

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))
		im := imdraw.New(nil)

		// Dibuja espacios de estacionamiento
		for i := 0; i < espacios; i++ {
			posX := win.Bounds().Center().X + tamanoEspacio*float64(i-espacios/2)
			im.Color = pixel.RGB(1, 1, 1)
			im.Push(pixel.V(posX, win.Bounds().Center().Y-tamanoEspacio/2))
			im.Push(pixel.V(posX+tamanoEspacio, win.Bounds().Center().Y+tamanoEspacio/2))
			im.Rectangle(0)

			if autos[i] {
				// Dibuja auto
				im.Color = pixel.RGB(0, 0, 1)
				offset := (tamanoEspacio - tamanoAuto) / 2
				im.Push(pixel.V(posX+offset, win.Bounds().Center().Y-tamanoAuto/2+offset))
				im.Push(pixel.V(posX+tamanoAuto+offset, win.Bounds().Center().Y+tamanoAuto/2+offset))
				im.Rectangle(0)
			}
		}

		im.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
