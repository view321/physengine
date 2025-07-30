package components

import "github.com/yohamta/donburi"

type MassData struct {
	Mass    float64
	InverseMass float64
	Inertia float64
	InverseInertia float64
}

var MassComponent = donburi.NewComponentType[MassData]()
