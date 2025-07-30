package scenes

import (
	"physengine/factory"
	Vec2 "physengine/helpers/vec2"
	"physengine/systems"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type MyScene struct {
	ecs  *ecs.ECS
	once sync.Once
}

func (ms *MyScene) Update() {
	ms.once.Do(ms.configure)
	ms.ecs.Update()
}

func (ms *MyScene) Draw(screen *ebiten.Image) {
	ms.ecs.Draw(screen)
}

func (ms *MyScene) configure() {
	ms.ecs = ecs.NewECS(donburi.NewWorld())
	ms.ecs.AddSystem(systems.UpdateCamera)
	ms.ecs.AddSystem(systems.UpdateCollisions)
	ms.ecs.AddSystem(systems.UpdateVelocity)
	ms.ecs.AddSystem(systems.UpdateTorque)
	ms.ecs.AddSystem(systems.UpdateAngularVelocity)
	ms.ecs.AddSystem(systems.UpdateForce)
	ms.ecs.AddSystem(systems.ApplyContinuousForce)
	ms.ecs.AddRenderer(0, systems.DrawCamera)
	factory.CreateCamera(ms.ecs)
	factory.CreateCollisionResolver(ms.ecs)

	// Create demo objects for polygon collision testing
	factory.CreatePolygonCollisionDemo(ms.ecs)

	// Create demo objects for rotation-aware collision testing
	factory.CreateRotatingCollisionDemo(ms.ecs)

	// Original demo objects
	factory.CreateTestSquare(ms.ecs, Vec2.Vec2{X: 0, Y: 400}, Vec2.Vec2{X: 0, Y: -50})
	factory.CreateTestCircle(ms.ecs, Vec2.Vec2{X: 100, Y: 400}, Vec2.Vec2{X: 0, Y: -50})
	factory.CreateRotatingObject(ms.ecs, Vec2.Vec2{X: -200, Y: 400}, Vec2.Vec2{X: 50, Y: -50}, 5000.0)
	factory.CreatePlatform(ms.ecs)
	factory.CreateConstForces(ms.ecs)
}
