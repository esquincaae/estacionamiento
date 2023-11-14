package models

import (
	"github.com/faiface/pixel"
	"math/rand"
	"time"
)

type Vehiculo struct {
	ID             			int
	Posicion       			pixel.Vec
	PosicionAnterior 		pixel.Vec
	Carril         			int
	Estacionado    			bool
	HoraSalida     			time.Time
	Entrando       			bool
	Teletransportando   	bool
	InicioTeletransporte 	time.Time
}

var CanalVehiculos chan Vehiculo

func Inicializar() {
	CanalVehiculos = make(chan Vehiculo)
	go GeneradorVehiculos()
}

func CrearVehiculo(id int) Vehiculo {
	MutexVehiculos.Lock()
	defer MutexVehiculos.Unlock()
	Vehiculo := Vehiculo{
		ID:       		id,
		Posicion: 		pixel.V(0, 300),
		Carril:   		-1,
		Estacionado: 	false,
	}
	Vehiculos = append(Vehiculos, Vehiculo)
	return Vehiculo
}

func EstablecerHoraSalida(vehiculo *Vehiculo) {
	rand.Seed(time.Now().UnixNano())
	salidaEn := time.Duration(rand.Intn(5)+1) * time.Second
	vehiculo.HoraSalida = time.Now().Add(salidaEn)
}

func ObtenerVehiculos() []Vehiculo {
	return Vehiculos
}

func AsignarCarrilAVehiculo(id int, carril int) {
	MutexVehiculos.Lock()
	defer MutexVehiculos.Unlock()
	for i := range Vehiculos {
		if Vehiculos[i].ID == id {
			Vehiculos[i].Carril = carril
		}
	}
}

func ReiniciarPosicionVehiculo(id int) {
	MutexVehiculos.Lock()
	defer MutexVehiculos.Unlock()
	for i := range Vehiculos {
		if Vehiculos[i].ID == id {
			Vehiculos[i].Posicion = pixel.V(0, 300)
		}
	}
}

func EncontrarPosicionVehiculo(id int) pixel.Vec {
	MutexVehiculos.Lock()
	defer MutexVehiculos.Unlock()
	for _, vehiculo := range Vehiculos {
		if vehiculo.ID == id {
			return vehiculo.Posicion
		}
	}
	return pixel.Vec{}
}

func EstacionarVehiculo(vehiculo *Vehiculo, destinoX, destinoY float64) {
	vehiculo.Posicion.X = destinoX
	vehiculo.Posicion.Y = destinoY
	vehiculo.Estacionado = true
	EstablecerHoraSalida(vehiculo)
}

func EliminarVehiculo(indice int) {
	Vehiculos = append(Vehiculos[:indice], Vehiculos[indice+1:]...)
}

func GeneradorVehiculos() {
	id := 0
	for {
		id++
		vehiculo := CrearVehiculo(id)
		CanalVehiculos <- vehiculo
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+500))
	}
}
