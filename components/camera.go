package components

import (
	"github.com/yohamta/donburi"
)

type CameraData struct {
	ViewportSizeX float64
	ViewportSizeY float64
}

var Camera = donburi.NewComponentType[CameraData]()
