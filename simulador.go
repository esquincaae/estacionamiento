package main

import (
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
)

const (
	anchoVentana         = 800.0
	altoVentana          = 600.0
	espacios             = 10 // 5 arriba y 5 abajo
	tamanoEspacio        = anchoVentana / float64(espacios/2+1)
	altoEspacio          = 100.0
	anchoAuto            = 40.0
	altoAuto             = 60.0
	velocidad            = 2.0
)

type Estacionamiento struct {
	espacios   int
	mu         sync.Mutex
}

type Auto struct {
	posX float64
	posY float64
	dir  float64 // -1 para los autos que salen, 1 para los autos que entran
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
	autos := make([]*Auto, 0)

	go func() {
		for {
			pos := e.Entrar()
			if pos != -1 {
				auto := &Auto{posX: 0 - anchoAuto, posY: altoVentana / 2 + altoAuto/2, dir: 1}
				autos = append(autos, auto)
				time.Sleep(time.Millisecond * 800)
			}

			if len(autos) > 0 && autos[0].posX > anchoVentana {
				e.Salir(0)
				autos = autos[1:]
			}
		}
	}()

	go func() {
		for {
			pos := e.Entrar()
			if pos != -1 {
				auto := &Auto{posX: anchoVentana, posY: altoVentana/2 - 3*altoAuto/2, dir: -1}
				autos = append(autos, auto)
				time.Sleep(time.Millisecond * 1300)
			}

			if len(autos) > 0 && autos[0].posX < 0-anchoAuto {
				e.Salir(0)
				autos = autos[1:]
			}
		}
	}()

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))
		im := imdraw.New(nil)

		// Dibuja estacionamiento
		im.Color = pixel.RGB(0.5, 0.5, 0.5)
		im.Push(pixel.V(0, altoVentana/2+altoEspacio))
		im.Push(pixel.V(anchoVentana, altoVentana/2-altoEspacio))
		im.Rectangle(0)

		// Dibuja espacios de estacionamiento
		for i := 0; i < espacios/2; i++ {
			posX := tamanoEspacio * float64(i+1)

			// Arriba
			im.Color = pixel.RGB(1, 1, 1)
			im.Push(pixel.V(posX, altoVentana/2+altoEspacio))
			im.Push(pixel.V(posX+tamanoEspacio, altoVentana))
			im.Rectangle(0)

			// Abajo
			im.Push(pixel.V(posX, 0))
			im.Push(pixel.V(posX+tamanoEspacio, altoVentana/2-altoEspacio))
			im.Rectangle(0)
		}

		// Dibuja autos
		for _, auto := range autos {
			im.Color = pixel.RGB(0, 0, 1)
			im.Push(pixel.V(auto.posX, auto.posY))
			im.Push(pixel.V(auto.posX+anchoAuto, auto.posY+altoAuto))
			im.Rectangle(0)

			auto.posX += velocidad * auto.dir
		}

		im.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
