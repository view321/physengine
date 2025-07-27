package components

import (
	"github.com/yohamta/donburi"
)

type DraggableData struct{
	IsDragged bool
}

var Draggable = donburi.NewComponentType[DraggableData]()