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
	ms.ecs.AddRenderer(0, systems.DrawCamera)
	factory.CreateCamera(ms.ecs)
	factory.CreateCollisionResolver(ms.ecs)
	factory.CreateTestCircle(ms.ecs, Vec2.Vec2{X: 0, Y: 300}, Vec2.Vec2{X: 0, Y: -50})
	factory.CreateTestCircle(ms.ecs, Vec2.Vec2{X: 0, Y: -300}, Vec2.Vec2{X: 0, Y: 50})
}
