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
	entity := ecs.World.Create(components.Transform, components.Player, components.Velocity, components.Drawable, components.CircleCollider)
	entry := ecs.World.Entry(entity)

	img, err := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	if err != nil {
		log.Printf("Failed to load player image: %v", err)
		// Return the entity without a sprite for now
		return entry
	}

	components.Drawable.SetValue(entry, components.DrawableData{Sprite: img})

	// Initialize transform with proper values
	tr := components.Transform.Get(entry)
	tr.Pos = Vec2.Vec2{X: 0, Y: 0} // Center of the screen
	tr.Scale = Vec2.Vec2{X: 1, Y: 1}
	tr.Rot = 0

	crcl := components.CircleCollider.Get(entry)
	crcl.Radius = 100

	//entity2 := ecs.World.Create(components.Transform, components.Drawable)
	//entry2 := ecs.World.Entry(entity2)
	//tr2 := components.Transform.Get(entry2)
	//tr2.Scale = Vec2.Vec2{X: 1, Y: 1}
	//tr2.Pos.X = 300
	//tr2.Pos.Y = 200
	//components.Drawable.SetValue(entry2, components.DrawableData{Sprite: img})
	//tr.Children = append(tr.Children, tr2)

	return entry
}
