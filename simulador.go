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
            columna := i % (espacios / 2)
            fila := i / (espacios / 2)
            auto.metaX = tamanoEspacio/2 + tamanoEspacio*float64(columna) + (anchoVentana-espacios/2*tamanoEspacio)/2 // centrar los espacios en el cuadro
            if fila == 0 {
                auto.metaY = altoVentana - altoEspacio - altoAuto/2 // fila superior dentro del cuadro
            } else {
                auto.metaY = altoEspacio + altoAuto/2 // fila inferior dentro del cuadro
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

	// Iniciar goroutines para simular los autos
	for i := 0; i < 5; i++ { // Simular 5 autos como ejemplo
		auto := &Auto{} // Crear una nueva instancia de Auto
		go e.simularAuto(auto) // Lanzar la simulación en una nueva goroutine
		time.Sleep(time.Second) // Esperar un segundo antes de lanzar la siguiente goroutine
	}

	for !win.Closed() {
        win.Clear(pixel.RGB(0.9, 0.9, 0.9)) // Establece el color de fondo a blanco grisáceo

		// Dibujar el borde del estacionamiento
		imd.Clear()
		imd.Color = pixel.RGB(0, 0, 0) // Establece el color a negro para el borde
		// Dibuja un rectángulo que representa el borde del estacionamiento
		imd.Push(pixel.V(tamanoEspacio, altoEspacio)) // Esquina inferior izquierda del estacionamiento
		imd.Push(pixel.V(anchoVentana-tamanoEspacio, altoVentana-altoEspacio)) // Esquina superior derecha del estacionamiento
		imd.Rectangle(3)

		// Dibujar la entrada/salida
		imd.Color = pixel.RGB(0.9, 0.9, 0.9) // Establece el color a blanco grisáceo para la entrada/salida
		imd.Push(pixel.V(tamanoEspacio, entradaY+altoAuto/2)) // Parte superior de la entrada/salida
		imd.Push(pixel.V(tamanoEspacio*1.5, entradaY-altoAuto/2)) // Parte inferior de la entrada/salida
		imd.Rectangle(3)

		 // Dibujar los espacios de estacionamiento
		 imd.Color = pixel.RGB(0.5, 0.5, 0.5) // Establece el color para los espacios de estacionamiento
		 for i := 0; i < espacios; i++ {
			 columna := i % (espacios / 2)
			 fila := i / (espacios / 2)
			 x := tamanoEspacio/2 + tamanoEspacio*float64(columna) + (anchoVentana-espacios/2*tamanoEspacio)/2 // centrar los espacios en el cuadro
			 var y float64
			 if fila == 0 {
				 y = altoVentana - altoEspacio - altoAuto/2 // fila superior dentro del cuadro
			 } else {
				 y = altoEspacio + altoAuto/2 // fila inferior dentro del cuadro
			 }
			 imd.Push(pixel.V(x-anchoAuto/2, y-altoAuto/2))
			 imd.Push(pixel.V(x+anchoAuto/2, y+altoAuto/2))
			 imd.Rectangle(1)
		 }

		// Dibujar los autos (esto es solo un ejemplo, necesitas tu propia lógica aquí)
		e.mu.Lock()
		for _, auto := range e.autos {
			imd.Color = pixel.RGB(0, 0, 1) // Establece el color a azul para los autos
			imd.Push(pixel.V(auto.posX-anchoAuto/2, auto.posY-altoAuto/2))
			imd.Push(pixel.V(auto.posX+anchoAuto/2, auto.posY+altoAuto/2))
			imd.Rectangle(0)
		}
		e.mu.Unlock()

		imd.Draw(win) // Dibuja todo lo que se ha añadido a imd en la ventana

		win.Update() // Actualiza la ventana para mostrar los cambios
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}
