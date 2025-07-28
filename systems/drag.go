package systems

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateDrag(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Draggable, components.CircleCollider, components.Transform))
	cam, _ := components.Camera.First(e.World)
	cam_comp := components.Camera.Get(cam)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		// Use the current mouse world position for collision detection
		mouse_world_x := cam_comp.LastMousePos.X
		mouse_world_y := cam_comp.LastMousePos.Y

		for entry := range query.Iter(e.World) {
			if components.PosInsideCollider(entry, Vec2.Vec2{X: mouse_world_x, Y: mouse_world_y}) {
				tr := components.Transform.Get(entry)
				// Calculate offset from mouse to object center
				mouse_offset := Vec2.Vec2{
					X: mouse_world_x - tr.Pos.X,
					Y: mouse_world_y - tr.Pos.Y,
				}
				components.Draggable.SetValue(entry, components.DraggableData{
					IsDragged:   true,
					MouseOffset: mouse_offset,
				})
			}
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		query.Each(e.World, func(entry *donburi.Entry) {
			components.Draggable.Get(entry).IsDragged = false
		})
	}

	// Update positions of dragged objects using mouse offset
	for entry := range query.Iter(e.World) {
		drag_comp := components.Draggable.Get(entry)
		if drag_comp.IsDragged {
			// Position object at mouse position minus the stored offset
			new_pos := Vec2.Vec2{
				X: cam_comp.LastMousePos.X - drag_comp.MouseOffset.X,
				Y: cam_comp.LastMousePos.Y - drag_comp.MouseOffset.Y,
			}
			components.SetPos(entry, new_pos)
		}
	}
}
