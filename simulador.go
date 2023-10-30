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

func (e *Estacionamiento) Entrar() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.espacios > 0 {
		e.espacios--
		return true
	}
	return false
}

func (e *Estacionamiento) Salir() {
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

	// Simulación de vehículos entrando y saliendo
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 500) // Nuevo vehículo cada 500ms
			if e.Entrar() {
				fmt.Println("Vehículo", i, "entró")
				time.Sleep(time.Second * 2) // El vehículo se queda por 2 segundos
				e.Salir()
				fmt.Println("Vehículo", i, "salió")
			} else {
				fmt.Println("Vehículo", i, "no encontró espacio")
			}
		}
	}()

	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		im := imdraw.New(nil)
		for i := 0; i < espacios; i++ {
			color := pixel.RGB(0, 1, 0) // Espacio libre (verde)
			if i >= espacios-e.espacios {
				color = pixel.RGB(1, 0, 0) // Espacio ocupado (rojo)
			}
			posX := win.Bounds().Center().X + tamanoEspacio*float64(i-espacios/2)
			im.Color = color
			im.Push(pixel.V(posX, win.Bounds().Center().Y-tamanoEspacio/2))
			im.Push(pixel.V(posX+tamanoEspacio, win.Bounds().Center().Y+tamanoEspacio/2))
			im.Rectangle(0)
		}
		im.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
