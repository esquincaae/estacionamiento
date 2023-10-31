package main

import (
	"math/rand"
	"sync"
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
)

const (
	anchoVentana   = 1024.0
	altoVentana    = 768.0
	espacios       = 20
	tamanoEspacio  = 40.0
	altoEspacio    = 100.0
	anchoAuto      = 30.0
	altoAuto       = 60.0
	velocidad      = 2.0
	entradaY       = altoVentana / 2
	distanciaEntrada = 300.0
)

type Auto struct {
	posX, posY, dirX, metaY float64
}

type Estacionamiento struct {
	espacios    [espacios]bool
	enEspera    []*Auto
	mu          sync.Mutex
}

func NuevoEstacionamiento() *Estacionamiento {
	return &Estacionamiento{}
}

func (e *Estacionamiento) estacionar() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, ocupado := range e.espacios {
		if !ocupado {
			e.espacios[i] = true
			return i
		}
	}
	return -1
}

func (e *Estacionamiento) salir(pos int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.espacios[pos] = false
}

func (e *Estacionamiento) simularAuto(auto *Auto) {
	// Anima el ingreso del auto al estacionamiento
	for auto.posX < anchoVentana-distanciaEntrada {
		auto.posX += velocidad
		time.Sleep(time.Millisecond * 10)
	}
	auto.dirX = 1
	auto.metaY = altoVentana - altoAuto/2

	pos := e.estacionar()
	if pos != -1 {
		targetX := tamanoEspacio*float64(pos%(espacios/2)) + tamanoEspacio/2
		if pos >= espacios/2 {
			auto.metaY = altoAuto/2
		}

		for auto.posX != targetX || auto.posY != auto.metaY {
			if auto.posX < targetX {
				auto.posX += velocidad
			} else if auto.posX > targetX {
				auto.posX -= velocidad
			}
			if auto.posY < auto.metaY {
				auto.posY += velocidad
			} else if auto.posY > auto.metaY {
				auto.posY -= velocidad
			}
			time.Sleep(time.Millisecond * 10)
		}
		auto.dirX = 0

		// Espera un tiempo aleatorio
		duracion := time.Second * time.Duration(rand.Intn(50)+10)
		time.Sleep(duracion)
		e.salir(pos)

		// Anima la salida del auto
		for auto.posX > 0 {
			auto.posX -= velocidad
			time.Sleep(time.Millisecond * 10)
		}
	} else {
		// Si no hay espacio, el auto se va
		for auto.posX > 0 {
			auto.posX -= velocidad
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Estacionamiento",
		Bounds: pixel.R(0, 0, anchoVentana, altoVentana),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	e := NuevoEstacionamiento()

	go func() {
		for {
			auto := &Auto{
				posX: 0,
				posY: entradaY,
				dirX: 1,
			}
			e.mu.Lock()
			e.enEspera = append(e.enEspera, auto)
			e.mu.Unlock()
			go e.simularAuto(auto)
			time.Sleep(time.Second * 2)
		}
	}()

	for !win.Closed() {
		win.Clear(pixel.RGB(0.9, 0.9, 0.9))

		// Dibujar estacionamiento
		imd := imdraw.New(nil)
		imd.Color = pixel.RGB(0.5, 0.5, 0.5)
		for i := 0; i < espacios/2; i++ {
			// Parte superior
			imd.Push(pixel.V(tamanoEspacio*float64(i), altoVentana-altoEspacio))
			imd.Push(pixel.V(tamanoEspacio*float64(i+1), altoVentana))
			imd.Rectangle(0)

			// Parte inferior
			imd.Push(pixel.V(tamanoEspacio*float64(i), 0))
			imd.Push(pixel.V(tamanoEspacio*float64(i+1), altoEspacio))
			imd.Rectangle(0)
		}
		imd.Draw(win)

		// Dibujar autos
		e.mu.Lock()
		for _, auto := range e.enEspera {
			imd := imdraw.New(nil)
			imd.Color = pixel.RGB(0, 0, 1)
			imd.Push(pixel.V(auto.posX, auto.posY))
			imd.Push(pixel.V(auto.posX+anchoAuto, auto.posY+altoAuto))
			imd.Rectangle(0)
			imd.Draw(win)
		}
		e.mu.Unlock()

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
