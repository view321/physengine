package factory

import (
	"log"
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlayer(ecs *ecs.ECS) *donburi.Entry {
	entity := ecs.World.Create(components.Transform, components.Player, components.Velocity, components.Drawable)
	entry := ecs.World.Entry(entity)

	img, err := assets.GetImage("D:/Coding/physengine/assets/assets/gopher.png")
	if err != nil {
		log.Printf("Failed to load player image: %v", err)
		// Return the entity without a sprite for now
		return entry
	}

	components.Drawable.SetValue(entry, components.DrawableData{Sprite: img})

	// Initialize transform with proper values
	tr := components.Transform.Get(entry)
	tr.Pos = Vec2.Vec2{X: 500, Y: 500} // Center of the screen
	tr.Scale = Vec2.Vec2{X: 1, Y: 1}
	tr.Rot = 0
	tr.ID = 1

	entity2 := ecs.World.Create(components.Transform, components.Drawable)
	entry2 := ecs.World.Entry(entity2)
	tr2 := components.Transform.Get(entry2)
	tr2.Scale = Vec2.Vec2{X: 1, Y: 1}
	tr2.Pos.X = 300
	tr2.Pos.Y = 200
	tr2.ID = 2
	components.Drawable.SetValue(entry2, components.DrawableData{Sprite: img})

	tr.Children = append(tr.Children, tr2)

	return entry
}
