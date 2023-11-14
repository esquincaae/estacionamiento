package models

import (
	"sync"
)

const (
	numCarriles  = 20
	AnchoCarril = 150.0
)

var (
	EstadoCarril [numCarriles]bool
	Vehiculos    []Vehiculo
	MutexVehiculos  sync.Mutex
	MutexCarril  sync.Mutex
	VehiculoEntrandoOSaliendo bool
)
