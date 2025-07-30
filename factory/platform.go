package factory

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlatform(e *ecs.ECS) {
	entity := e.World.Create(components.Transform, components.MaterialComponent)
	entry := e.World.Entry(entity)

	// Position the platform at the center of where we want it
	platformCenter := Vec2.Vec2{X: 0, Y: -200}
	components.SetPos(entry, platformCenter)

	// Define AABB bounds relative to the object center
	// Platform should be 800 units wide and 100 units tall
	halfWidth := 600.0
	halfHeight := 50.0
	aabbMin := Vec2.Vec2{X: -halfWidth, Y: -halfHeight}
	aabbMax := Vec2.Vec2{X: halfWidth, Y: halfHeight}

	// Add AABB component with relative bounds
	donburi.Add(entry, components.AABB_Component, &components.AABB_Data{Min: aabbMin, Max: aabbMax})
	donburi.Add(entry, components.MassComponent, &components.MassData{
		Mass:           0,
		InverseMass:    0,
		Inertia:        0,
		InverseInertia: 0,
	})
	donburi.Add(entry, components.Velocity, &components.VelocityData{})
	donburi.Add(entry, components.AngularVelocity, &components.AngularVelocityData{})

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0
	mat.StaticFriction = 1.0  // Very high static friction
	mat.DynamicFriction = 0.8 // Very high dynamic friction
}
