package components

import "github.com/yohamta/donburi"

type MassData struct {
	Mass float64
}

var MassComponent = donburi.NewComponentType[MassData]()