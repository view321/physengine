package components

import "github.com/yohamta/donburi"

type MaterialData struct{
	Density float64
	Restitution float64
	StaticFriction float64
	DynamicFriction float64
}

var MaterialComponent = donburi.NewComponentType[MaterialData]()