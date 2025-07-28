package components

import (
	Vec2 "physengine/helpers/vec2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type CameraData struct {
	ViewportSizeX      float64
	ViewportSizeY      float64
	Zoom               Vec2.Vec2
	LastMousePos       Vec2.Vec2 // World coordinates for collision detection
	LastScreenMousePos Vec2.Vec2 // Screen coordinates for delta calculation
	MouseDelta         Vec2.Vec2
	Op                 ebiten.DrawImageOptions
}

var Camera = donburi.NewComponentType[CameraData]()

func ChangeZoom(w donburi.World, new_zoom Vec2.Vec2) {
	cam_entry, _ := Camera.First(w)
	old_cam_obj := Camera.Get(cam_entry)
	old_cam_obj.Zoom = new_zoom
}
