package scenes

import (
	"physengine/factory"
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
	ms.ecs.AddRenderer(0, systems.DrawCamera)
	factory.CreateCamera(ms.ecs)
}
