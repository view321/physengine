package factory

import (
	"fmt"
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlatform(e *ecs.ECS) {
	entity := e.World.Create(components.Transform, components.MaterialComponent, components.Drawable)
	entry := e.World.Entry(entity)
	// Position the platform in the middle of the screen where objects can land on it
	donburi.Add(entry, components.AABB_Component, &components.AABB_Data{Min: Vec2.Vec2{X: -400, Y: 100}, Max: Vec2.Vec2{X: 400, Y: 200}})
	donburi.Add(entry, components.MassComponent, &components.MassData{
		Mass:           0,
		InverseMass:    0,
		Inertia:        0,
		InverseInertia: 0,
	})
	donburi.Add(entry, components.Velocity, &components.VelocityData{})
	donburi.Add(entry, components.AngularVelocity, &components.AngularVelocityData{})

	// Add a visual representation for the platform
	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 8, Y: 1} // Make it wide and flat

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.5

	fmt.Printf("Platform created at position: %v, AABB: Min=%v, Max=%v\n",
		components.Transform.Get(entry).Pos,
		components.AABB_Component.Get(entry).Min,
		components.AABB_Component.Get(entry).Max)
}
