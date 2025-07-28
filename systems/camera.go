package systems

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateCamera(e *ecs.ECS) {
	cam_entry, _ := components.Camera.First(e.World)
	cam_comp := components.Camera.Get(cam_entry)
	cam_tr := components.Transform.Get(cam_entry)

	// Get current mouse position in screen coordinates
	current_mouse_pos_x, current_mouse_pos_y := ebiten.CursorPosition()
	current_screen_pos := Vec2.Vec2{X: float64(current_mouse_pos_x), Y: float64(current_mouse_pos_y)}

	// Calculate mouse delta in screen coordinates (more responsive)
	screen_delta_x := current_screen_pos.X - cam_comp.LastScreenMousePos.X
	screen_delta_y := current_screen_pos.Y - cam_comp.LastScreenMousePos.Y

	// Use screen delta directly for dragging (1:1 mapping)
	cam_comp.MouseDelta = Vec2.Vec2{X: screen_delta_x, Y: screen_delta_y}

	// Convert screen coordinates to world coordinates for collision detection
	// Invert Y axis so it points up
	mouse_world_x := (current_screen_pos.X-cam_comp.ViewportSizeX/2)/cam_comp.Zoom.X + cam_tr.Pos.X
	mouse_world_y := -(current_screen_pos.Y-cam_comp.ViewportSizeY/2)/cam_comp.Zoom.Y + cam_tr.Pos.Y

	// Update both screen and world mouse positions
	cam_comp.LastScreenMousePos = current_screen_pos
	cam_comp.LastMousePos = Vec2.Vec2{X: mouse_world_x, Y: mouse_world_y}
}

func DrawCamera(e *ecs.ECS, screen_camera *ebiten.Image) {
	camera, _ := components.Camera.First(e.World)
	camera_tr := components.Transform.Get(camera)
	camera_comp := components.Camera.Get(camera)
	query := donburi.NewQuery(filter.Contains(components.Transform, components.Drawable))

	for entry := range query.Iter(e.World) {
		obj_tr := components.Transform.Get(entry)
		obj_drawable := components.Drawable.Get(entry)

		// Create a new DrawImageOptions for each entity to avoid state issues
		op := &ebiten.DrawImageOptions{}

		// Calculate world position relative to camera
		world_x := obj_tr.Pos.X - camera_tr.Pos.X
		world_y := obj_tr.Pos.Y - camera_tr.Pos.Y

		// Convert world coordinates to screen coordinates
		// Invert Y axis so it points up
		screen_x := world_x*camera_comp.Zoom.X + camera_comp.ViewportSizeX/2
		screen_y := -world_y*camera_comp.Zoom.Y + camera_comp.ViewportSizeY/2

		// Apply transformations in the correct order: center, scale, rotate, translate
		// First center the sprite on its origin
		op.GeoM.Translate(-float64(obj_drawable.Sprite.Bounds().Dx())/2, -float64(obj_drawable.Sprite.Bounds().Dy())/2)
		// Then scale (including zoom)
		op.GeoM.Scale(obj_tr.Scale.X*camera_comp.Zoom.X, obj_tr.Scale.Y*camera_comp.Zoom.Y)
		// Then rotate around the center (invert rotation for Y-axis inversion)
		op.GeoM.Rotate(-obj_tr.Rot)
		// Finally translate to screen position
		op.GeoM.Translate(screen_x, screen_y)

		screen_camera.DrawImage(obj_drawable.Sprite, op)
	}
}
