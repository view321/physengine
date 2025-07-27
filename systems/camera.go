package systems

import (
	"physengine/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func DrawCamera(e *ecs.ECS, screen_camera *ebiten.Image) {
	camera, _ := components.Camera.First(e.World)
	camera_tr := components.Transform.Get(camera)
	camera_comp := components.Camera.Get(camera)
	query := donburi.NewQuery(filter.Contains(components.Transform, components.Drawable))

	for entry := range query.Iter(e.World) {
		obj_tr := components.Transform.Get(entry)
		obj_drawable := components.Drawable.Get(entry)

		// Create a new DrawImageOptions for each entity
		op := &ebiten.DrawImageOptions{}

		// Calculate world position relative to camera
		world_x := obj_tr.Pos.X - camera_tr.Pos.X
		world_y := obj_tr.Pos.Y - camera_tr.Pos.Y

		// Convert world coordinates to screen coordinates
		screen_x := world_x + camera_comp.ViewportSizeX/2
		screen_y := world_y + camera_comp.ViewportSizeY/2

		// Apply transformations in the correct order: center, scale, rotate, translate
		// First center the sprite on its origin
		op.GeoM.Translate(-float64(obj_drawable.Sprite.Bounds().Dx())/2, -float64(obj_drawable.Sprite.Bounds().Dy())/2)
		// Then scale
		op.GeoM.Scale(obj_tr.Scale.X, obj_tr.Scale.Y)
		// Then rotate around the center
		op.GeoM.Rotate(obj_tr.Rot)
		// Finally translate to screen position
		op.GeoM.Translate(screen_x, screen_y)

		screen_camera.DrawImage(obj_drawable.Sprite, op)
	}
}
