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
	cam_tr := components.Transform.Get(cam)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		mouse_screen_pos_x, mouse_screen_pos_y := ebiten.CursorPosition()
		mouse_world_x := mouse_screen_pos_x - int(cam_tr.Pos.X)
		mouse_world_y := mouse_screen_pos_y - int(cam_tr.Pos.Y)
		for entry := range query.Iter(e.World) {
			col := components.CircleCollider.Get(entry)
			if components.PosInsideCollider(col, &Vec2.Vec2{X: float64(mouse_world_x), Y: float64(mouse_world_y)}) {
				components.Draggable.SetValue(entry, components.DraggableData{IsDragged: true})
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		components.Draggable.Each(e.World, func(entry *donburi.Entry) {
			components.Draggable.Get(entry).IsDragged = false
		})
	}
	//Не закончено
}
