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
	anchoVentana    = 1024.0
	altoVentana     = 768.0
	espacios        = 20
	tamanoEspacio   = 40.0
	altoEspacio     = 100.0
	anchoAuto       = 30.0
	altoAuto        = 60.0
	velocidad       = 2.0
	entradaY        = altoVentana / 2
	distanciaEntrada = 300.0
	centroX         = anchoVentana / 2
	centroY         = altoVentana / 2
)

type Auto struct {
	posX, posY, dirX, dirY, metaX, metaY float64
	estacionado                          bool
	tiempoEstacionado                    time.Duration
}

type Estacionamiento struct {
	espacios [espacios]bool
	autos    []*Auto
	mu       sync.Mutex
}

func NuevoEstacionamiento() *Estacionamiento {
	return &Estacionamiento{}
}

func (e *Estacionamiento) estacionar(auto *Auto) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, ocupado := range e.espacios {
		if !ocupado {
			e.espacios[i] = true
			auto.metaX = tamanoEspacio*float64(i%(espacios/2)) + tamanoEspacio/2 // X es el centro de la columna
			fila := i / (espacios / 2)
			if fila == 0 {
				auto.metaY = altoVentana - altoEspacio/2 // Y es el centro para fila 0
			} else {
				auto.metaY = altoEspacio/2 // Y es el centro para fila 1
			}
			e.autos = append(e.autos, auto)
			return i
		}
	}
	return -1
}

func (e *Estacionamiento) salir(auto *Auto) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, a := range e.autos {
		if a == auto {
			e.espacios[i] = false
			e.autos = append(e.autos[:i], e.autos[i+1:]...)
			break
		}
	}
}

func (e *Estacionamiento) simularAuto(auto *Auto) {
	auto.posX = 0
	auto.posY = entradaY
	auto.dirX = 1

	for auto.posX < distanciaEntrada {
		auto.posX += velocidad
		time.Sleep(time.Millisecond * 10)
	}
	auto.dirX = 0

	for {
		auto.estacionado = false
		time.Sleep(time.Second * time.Duration(rand.Intn(10)+1)) // Esperar un momento antes de intentar estacionarse
		pos := e.estacionar(auto)
		if pos != -1 {
			// Mover al auto al centro de la pantalla
			for auto.posX != auto.metaX || auto.posY != auto.metaY {
				if auto.posX < auto.metaX {
					auto.posX += velocidad
				} else if auto.posX > auto.metaX {
					auto.posX -= velocidad
				}
				if auto.posY < auto.metaY {
					auto.posY += velocidad
				} else if auto.posY > auto.metaY {
					auto.posY -= velocidad
				}
				time.Sleep(time.Millisecond * 10)
			}

			auto.estacionado = true

			auto.tiempoEstacionado = time.Second * time.Duration(rand.Intn(50)+10)
			time.Sleep(auto.tiempoEstacionado)

			// Mover al auto a una posición vacía
			auto.dirX = -velocidad
			for auto.posX > 0 {
				auto.posX += auto.dirX
				time.Sleep(time.Millisecond * 10)
			}
			e.salir(auto)
		} else {
			// Si no hay espacio, el auto se va
			for auto.posX > 0 {
				auto.posX += auto.dirX
				time.Sleep(time.Millisecond * 10)
			}
		}
		time.Sleep(time.Second * time.Duration(rand.Intn(100)+10))
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

	imd := imdraw.New(nil) // Crear una sola instancia de imdraw.IMDraw

	go func() {
		for {
			auto := &Auto{}
			go e.simularAuto(auto)
			time.Sleep(time.Second * time.Duration(rand.Intn(10)+1))
		}
	}()

	for !win.Closed() {
		win.Clear(pixel.RGB(0.9, 0.9, 0.9))

		// Dibujar estacionamiento
		imd.Clear() // Limpiar la instancia imd antes de dibujar las divisiones
		imd.Color = pixel.RGB(0.5, 0.5, 0.5)
		for i := 0; i < espacios/2; i++ {
			// Parte superior
			imd.Push(pixel.V(tamanoEspacio*float64(i)+tamanoEspacio/2, altoVentana-altoEspacio/2))
			imd.Push(pixel.V(tamanoEspacio*float64(i+1)+tamanoEspacio/2, altoVentana))
			imd.Rectangle(0)

			// Parte inferior
			imd.Push(pixel.V(tamanoEspacio*float64(i)+tamanoEspacio/2, altoEspacio/2))
			imd.Push(pixel.V(tamanoEspacio*float64(i+1)+tamanoEspacio/2, 0))
			imd.Rectangle(0)
		}
		imd.Draw(win)

		// Dibujar autos
		e.mu.Lock()
		for _, auto := range e.autos {
			imd.Color = pixel.RGB(0, 0, 1)
			imd.Push(pixel.V(auto.posX-anchoAuto/2, auto.posY-altoAuto/2))
			imd.Push(pixel.V(auto.posX+anchoAuto/2, auto.posY+altoAuto/2))
			imd.Rectangle(0)
		}
		e.mu.Unlock()

		imd.Draw(win)

		win.Update()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}
