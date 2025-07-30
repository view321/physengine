package systems

import (
	"physengine/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateVelocity(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Velocity))
	for entry := range query.Iter(e.World) {
		tr := components.Transform.Get(entry)
		vel := components.Velocity.Get(entry)
		deltaTime := float64(e.Time.DeltaTime().Seconds())
		newPos := tr.Pos.Add(vel.Velocity.Mult(deltaTime))
		components.SetPos(entry, newPos)
	}
}
