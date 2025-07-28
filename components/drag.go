package components

import (
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type DraggableData struct {
	IsDragged   bool
	MouseOffset Vec2.Vec2 // Offset from mouse to object center when dragging started
}

var Draggable = donburi.NewComponentType[DraggableData]()
