package scenes

import (
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"carro/models"
	"carro/views"
)

func Run() {

	models.Inicializar()

	ventana, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Simulación de Estacionamiento",
		Bounds: pixel.R(0, 0, 800, 600),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for vehiculo := range models.CanalVehiculos {
			models.MutexCarril.Lock()
			for _, ocupado := range models.EstadoCarril {
				if !ocupado {
					break
				}
			}
			models.MutexCarril.Unlock()

			go models.Carril(vehiculo.ID)
		}
	}()

	for !ventana.Closed() {
		ventana.Clear(colornames.White)
		views.DibujarEstacionamiento(ventana, models.ObtenerVehiculos())
		ventana.Update()
		models.MutexVehiculos.Lock()
		models.LogicaMovimientoVehiculos()
		models.MutexVehiculos.Unlock()

		time.Sleep(16 * time.Millisecond)
	}
}
