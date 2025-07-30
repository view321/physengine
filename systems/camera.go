package systems

import (
	"image/color"
	"math"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	query2 := donburi.NewQuery(filter.Contains(components.AABB_Component))
	for entry := range query2.Iter(e.World) {
		aabb := components.AABB_Component.Get(entry)
		obj_tr := components.Transform.Get(entry)

		// AABB bounds are in world coordinates relative to object position
		// Calculate the four corners of the AABB in world space
		worldCorners := []Vec2.Vec2{
			{X: obj_tr.Pos.X + aabb.Min.X, Y: obj_tr.Pos.Y + aabb.Min.Y}, // bottom-left
			{X: obj_tr.Pos.X + aabb.Max.X, Y: obj_tr.Pos.Y + aabb.Min.Y}, // bottom-right
			{X: obj_tr.Pos.X + aabb.Max.X, Y: obj_tr.Pos.Y + aabb.Max.Y}, // top-right
			{X: obj_tr.Pos.X + aabb.Min.X, Y: obj_tr.Pos.Y + aabb.Max.Y}, // top-left
		}

		// Apply rotation around object center if object is rotated
		if obj_tr.Rot != 0 {
			objectCenter := obj_tr.Pos
			for i := range worldCorners {
				// Translate to origin
				relativeX := worldCorners[i].X - objectCenter.X
				relativeY := worldCorners[i].Y - objectCenter.Y

				// Apply rotation
				rotatedX := relativeX*math.Cos(obj_tr.Rot) - relativeY*math.Sin(obj_tr.Rot)
				rotatedY := relativeX*math.Sin(obj_tr.Rot) + relativeY*math.Cos(obj_tr.Rot)

				// Translate back
				worldCorners[i].X = objectCenter.X + rotatedX
				worldCorners[i].Y = objectCenter.Y + rotatedY
			}
		}

		// Convert world coordinates to screen coordinates
		screenCorners := make([]Vec2.Vec2, 4)
		for i, worldCorner := range worldCorners {
			// Calculate relative to camera
			relativeX := worldCorner.X - camera_tr.Pos.X
			relativeY := worldCorner.Y - camera_tr.Pos.Y

			// Convert to screen coordinates (invert Y axis)
			screenCorners[i] = Vec2.Vec2{
				X: relativeX*camera_comp.Zoom.X + camera_comp.ViewportSizeX/2,
				Y: -relativeY*camera_comp.Zoom.Y + camera_comp.ViewportSizeY/2,
			}
		}

		// Draw the AABB wireframe
		vector.StrokeLine(screen_camera, float32(screenCorners[0].X), float32(screenCorners[0].Y), float32(screenCorners[1].X), float32(screenCorners[1].Y), 2, color.White, false)
		vector.StrokeLine(screen_camera, float32(screenCorners[1].X), float32(screenCorners[1].Y), float32(screenCorners[2].X), float32(screenCorners[2].Y), 2, color.White, false)
		vector.StrokeLine(screen_camera, float32(screenCorners[2].X), float32(screenCorners[2].Y), float32(screenCorners[3].X), float32(screenCorners[3].Y), 2, color.White, false)
		vector.StrokeLine(screen_camera, float32(screenCorners[3].X), float32(screenCorners[3].Y), float32(screenCorners[0].X), float32(screenCorners[0].Y), 2, color.White, false)
	}
	query3 := donburi.NewQuery(filter.Contains(components.CircleCollider))
	for entry := range query3.Iter(e.World) {
		crcl := components.CircleCollider.Get(entry)
		obj_tr := components.Transform.Get(entry)

		// Calculate world position relative to camera
		world_x := obj_tr.Pos.X - camera_tr.Pos.X
		world_y := obj_tr.Pos.Y - camera_tr.Pos.Y

		// Convert world coordinates to screen coordinates
		// Invert Y axis so it points up (consistent with sprites and AABB)
		screen_x := world_x*camera_comp.Zoom.X + camera_comp.ViewportSizeX/2
		screen_y := -world_y*camera_comp.Zoom.Y + camera_comp.ViewportSizeY/2

		// Scale the radius by the camera zoom (use average of X and Y zoom for consistency)
		radius_scale := (camera_comp.Zoom.X + camera_comp.Zoom.Y) / 2
		scaled_radius := crcl.Radius * radius_scale

		vector.StrokeCircle(screen_camera, float32(screen_x), float32(screen_y), float32(scaled_radius), 2, color.White, false)
	}

	// Draw polygon wireframes
	query4 := donburi.NewQuery(filter.Contains(components.PolygonCollider))
	for entry := range query4.Iter(e.World) {
		// Get world vertices (already transformed with rotation and position)
		worldVertices := components.GetWorldVertices(entry)
		if worldVertices == nil || len(worldVertices) < 3 {
			continue
		}

		// Convert vertices to screen coordinates
		screenVertices := make([]Vec2.Vec2, len(worldVertices))
		for i, vertex := range worldVertices {
			// Calculate world position relative to camera
			world_x := vertex.X - camera_tr.Pos.X
			world_y := vertex.Y - camera_tr.Pos.Y

			// Convert world coordinates to screen coordinates
			// Invert Y axis so it points up
			screen_x := world_x*camera_comp.Zoom.X + camera_comp.ViewportSizeX/2
			screen_y := -world_y*camera_comp.Zoom.Y + camera_comp.ViewportSizeY/2

			screenVertices[i] = Vec2.Vec2{X: screen_x, Y: screen_y}
		}

		// Choose color based on polygon type
		var wireframeColor color.Color
		switch len(worldVertices) {
		case 3:
			wireframeColor = color.RGBA{255, 0, 0, 255} // Red for triangles
		case 4:
			wireframeColor = color.RGBA{0, 255, 0, 255} // Green for rectangles
		case 5:
			wireframeColor = color.RGBA{0, 0, 255, 255} // Blue for pentagons
		case 6:
			wireframeColor = color.RGBA{255, 255, 0, 255} // Yellow for hexagons
		default:
			wireframeColor = color.RGBA{255, 0, 255, 255} // Magenta for other polygons
		}

		// Draw polygon wireframe by connecting vertices
		for i := 0; i < len(screenVertices); i++ {
			next := (i + 1) % len(screenVertices)
			current := screenVertices[i]
			nextVertex := screenVertices[next]

			vector.StrokeLine(
				screen_camera,
				float32(current.X), float32(current.Y),
				float32(nextVertex.X), float32(nextVertex.Y),
				2, wireframeColor, false,
			)
		}

		// Draw vertex markers
		for _, vertex := range screenVertices {
			vector.StrokeCircle(
				screen_camera,
				float32(vertex.X), float32(vertex.Y),
				3, 1, color.White, false,
			)
		}
	}
}
func ApplyRotToPoint(p1 *Vec2.Vec2, rot float64) {
	oldX := p1.X
	p1.X = p1.X*math.Cos(rot) - p1.Y*math.Sin(rot)
	p1.Y = oldX*math.Sin(rot) + p1.Y*math.Cos(rot)
}

func ApplyRotToPointAroundCenter(p1 *Vec2.Vec2, center Vec2.Vec2, rot float64) {
	// Translate point relative to center
	relativeX := p1.X - center.X
	relativeY := p1.Y - center.Y

	// Apply rotation
	oldX := relativeX
	relativeX = relativeX*math.Cos(rot) - relativeY*math.Sin(rot)
	relativeY = oldX*math.Sin(rot) + relativeY*math.Cos(rot)

	// Translate back
	p1.X = center.X + relativeX
	p1.Y = center.Y + relativeY
}
