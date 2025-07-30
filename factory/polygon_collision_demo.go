package factory

import (
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// CreatePolygonCollisionDemo creates objects to demonstrate polygon collision system
func CreatePolygonCollisionDemo(ecs *ecs.ECS) {
	// Create different polygon shapes
	CreateTestTriangle(ecs, Vec2.Vec2{X: 200, Y: 150}, Vec2.Vec2{X: 80, Y: 20})
	CreateTestPentagon(ecs, Vec2.Vec2{X: 400, Y: 200}, Vec2.Vec2{X: -60, Y: 40})
	CreateTestHexagon(ecs, Vec2.Vec2{X: 600, Y: 250}, Vec2.Vec2{X: -40, Y: -30})
	CreateTestRectangle(ecs, Vec2.Vec2{X: 300, Y: 400}, Vec2.Vec2{X: 50, Y: -50}, 120, 80)

	// Create some circles to interact with polygons
	CreatePolygonDemoCircle(ecs, Vec2.Vec2{X: 100, Y: 300}, Vec2.Vec2{X: 100, Y: -20})
	CreatePolygonDemoCircle(ecs, Vec2.Vec2{X: 500, Y: 100}, Vec2.Vec2{X: -30, Y: 80})

	// Create some AABB boxes to interact with polygons
	CreatePolygonDemoBox(ecs, Vec2.Vec2{X: 700, Y: 200}, Vec2.Vec2{X: -80, Y: 30})
	CreatePolygonDemoBox(ecs, Vec2.Vec2{X: 150, Y: 450}, Vec2.Vec2{X: 40, Y: -60})

	// Create stationary objects
	CreateStationaryPolygon(ecs, Vec2.Vec2{X: 800, Y: 400}, Vec2.Vec2{X: 0, Y: 0})
	CreateStationaryCircle(ecs, Vec2.Vec2{X: 50, Y: 50}, Vec2.Vec2{X: 0, Y: 0})
}

// CreatePolygonDemoCircle creates a circle for polygon collision demo
func CreatePolygonDemoCircle(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.CircleCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.Torque)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 1.5)

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/enemy.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 0.8, Y: 0.8}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 7
	mc.InverseMass = 1 / mc.Mass

	crcl := components.CircleCollider.Get(entry)
	crcl.Radius = 45
	mc.Inertia = mc.Mass * crcl.Radius * crcl.Radius * 0.5
	mc.InverseInertia = 1 / mc.Inertia

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.75
	mat.StaticFriction = 0.35
	mat.DynamicFriction = 0.25

	return entry
}

// CreatePolygonDemoBox creates an AABB box for polygon collision demo
func CreatePolygonDemoBox(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.AABB_Component, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.Torque)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 2.0)

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 0.9, Y: 0.9}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 9
	mc.InverseMass = 1 / mc.Mass
	mc.Inertia = mc.Mass * 9000
	mc.InverseInertia = 1 / mc.Inertia

	box := components.AABB_Component.Get(entry)
	box.Min = Vec2.Vec2{X: -50, Y: -40}
	box.Max = Vec2.Vec2{X: 50, Y: 40}

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.7
	mat.StaticFriction = 0.4
	mat.DynamicFriction = 0.3

	return entry
}

// CreateStationaryPolygon creates a stationary polygon for collision testing
func CreateStationaryPolygon(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.PolygonCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.Torque)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 0.0) // No rotation

	// Create a large octagon
	vertices := components.CreateRegularPolygon(8, 80)

	poly := components.PolygonCollider.Get(entry)
	poly.Vertices = vertices

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/player.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1.2, Y: 1.2}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 25 // Heavy object
	mc.InverseMass = 1 / mc.Mass
	mc.Inertia = mc.Mass * 20000
	mc.InverseInertia = 1 / mc.Inertia

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.5
	mat.StaticFriction = 0.6
	mat.DynamicFriction = 0.4

	return entry
}

// CreateStationaryCircle creates a stationary circle for collision testing
func CreateStationaryCircle(ecs *ecs.ECS, pos Vec2.Vec2, vel Vec2.Vec2) *donburi.Entry {
	entity := ecs.World.Create(components.MaterialComponent, components.Transform, components.CircleCollider, components.Drawable, components.MassComponent, components.Velocity, components.AngularVelocity, components.Torque)
	entry := ecs.World.Entry(entity)
	components.SetPos(entry, pos)
	components.Velocity.Get(entry).Velocity = vel
	components.SetAngularVelocity(entry, 0.0) // No rotation

	img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/enemy.png")
	components.Drawable.Get(entry).Sprite = img
	components.Transform.Get(entry).Scale = Vec2.Vec2{X: 1.1, Y: 1.1}

	mc := components.MassComponent.Get(entry)
	mc.Mass = 20 // Heavy object
	mc.InverseMass = 1 / mc.Mass

	crcl := components.CircleCollider.Get(entry)
	crcl.Radius = 90
	mc.Inertia = mc.Mass * crcl.Radius * crcl.Radius * 0.5
	mc.InverseInertia = 1 / mc.Inertia

	mat := components.MaterialComponent.Get(entry)
	mat.Restitution = 0.6
	mat.StaticFriction = 0.5
	mat.DynamicFriction = 0.35

	return entry
}
