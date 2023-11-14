package models

import (
	"math/rand"
	"time"
)

func EsperarPorPosicion(id int, destinoX float64) {
	for {
		posVehiculo := EncontrarPosicionVehiculo(id)
		if posVehiculo.X >= destinoX {
			break
		}
		time.Sleep(16 * time.Millisecond)
	}
}

func EncontrarCarrilDisponible() (int, bool) {
	MutexCarril.Lock()
	defer MutexCarril.Unlock()
	rand.Seed(time.Now().UnixNano())
	carriles := rand.Perm(numCarriles)
	for _, c := range carriles {
		if !EstadoCarril[c] {
			EstadoCarril[c] = true
			return c, true
		}
	}
	return -1, false
}

func Carril(id int) {
	CrearVehiculo(id)
	EsperarPorPosicion(id, 100)
	carril, carrilEncontrado := EncontrarCarrilDisponible()
	if !carrilEncontrado {
		ReiniciarPosicionVehiculo(id)
		return
	}
	AsignarCarrilAVehiculo(id, carril)
}
