package factory

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func CreateConstForces(e *ecs.ECS) {
	var entity donburi.Entity = e.World.Create(components.ContForceComp)
	entry := e.World.Entry(entity)
	comp := components.ContForceComp.Get(entry)
	gravity := components.ContForceContainer{}
	gravity.Force = Vec2.Vec2{X: 0, Y: -10000}
	query := donburi.NewQuery(filter.Contains(components.GravityTag))
	for entry := range query.Iter(e.World) {
		gravity.Objects = append(gravity.Objects, entry)
	}
	comp.Containers = append(comp.Containers, gravity)
}
