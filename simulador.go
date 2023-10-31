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
	anchoVentana         = 800.0
	altoVentana          = 600.0
	espacios             = 10
	tamanoEspacio        = anchoVentana / float64(espacios+1)
	altoEspacio          = 100.0
	anchoAuto            = 40.0
	altoAuto             = 80.0
	velocidad            = 1.0
	grosorLineaDivisoria = 2.0
	distanciaEntreAutos  = 10.0
)

type Estacionamiento struct {
	espacios   int
	mu         sync.Mutex
	ocupados   []*Auto
	enEspera   []*Auto
}

type Auto struct {
	posX float64
	posY float64
	dir  float64
}

func NuevoEstacionamiento(capacidad int) *Estacionamiento {
	return &Estacionamiento{
		espacios: capacidad,
		ocupados: make([]*Auto, capacidad),
	}
}

func (e *Estacionamiento) Entrar(auto *Auto) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, lugar := range e.ocupados {
		if lugar == nil {
			e.ocupados[i] = auto
			return i
		}
	}

	e.enEspera = append(e.enEspera, auto)
	return -1
}

func (e *Estacionamiento) Salir(i int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.ocupados[i] = nil
	if len(e.enEspera) > 0 {
		e.ocupados[i] = e.enEspera[0]
		e.enEspera = e.enEspera[1:]
	}
}

func run() {
	rand.Seed(time.Now().UnixNano())

	cfg := pixelgl.WindowConfig{
		Title:  "Estacionamiento",
		Bounds: pixel.R(0, 0, anchoVentana, altoVentana),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	e := NuevoEstacionamiento(espacios)

	go func() {
		for {
			auto := &Auto{posX: -anchoAuto - distanciaEntreAutos, posY: altoVentana/2 + altoAuto/2, dir: 1}
			pos := e.Entrar(auto)

			if pos != -1 {
				// Cuando un auto encuentra espacio, se estaciona y después de un tiempo se va
				go func(p int) {
					time.Sleep(time.Duration(rand.Intn(15)+5) * time.Second)
					e.Salir(p)
				}(pos)
			}
			time.Sleep(time.Millisecond * 1500)
		}
	}()

	for !win.Closed() {
		win.Clear(pixel.RGB(0.8, 0.8, 0.8)) // Background claro
		im := imdraw.New(nil)

		// Dibuja estacionamientos
		for i := 0; i < espacios/2; i++ {
			im.Color = pixel.RGB(0.5, 0.5, 0.5)
			xStart := tamanoEspacio * float64(i)
			im.Push(pixel.V(xStart, altoVentana/2+altoEspacio))
			im.Push(pixel.V(xStart+tamanoEspacio-grosorLineaDivisoria, altoVentana/2))
			im.Rectangle(0)

			im.Push(pixel.V(xStart, altoVentana/2-altoEspacio))
			im.Push(pixel.V(xStart+tamanoEspacio-grosorLineaDivisoria, altoVentana/2))
			im.Rectangle(0)
		}

		// Dibuja autos
		e.mu.Lock()
		for i, auto := range e.ocupados {
			if auto != nil {
				if i%2 == 0 {
					auto.posY = altoVentana/2 + altoAuto/2
				} else {
					auto.posY = altoVentana/2 - altoAuto/2 - altoEspacio
				}

				im.Color = pixel.RGB(0, 0, 1)
				im.Push(pixel.V(auto.posX, auto.posY))
				im.Push(pixel.V(auto.posX+anchoAuto, auto.posY+altoAuto))
				im.Rectangle(0)

				if auto.posX < tamanoEspacio*float64(i/2) || auto.dir == -1 {
					auto.posX += velocidad * auto.dir
				}
			}
		}
		e.mu.Unlock()

		im.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

