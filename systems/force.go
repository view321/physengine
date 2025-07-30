package systems

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateForce(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.ForceComponent, components.Velocity, components.MassComponent))
	for entry := range query.Iter(e.World) {
		frc := components.ForceComponent.Get(entry)
		mss := components.MassComponent.Get(entry)
		vel := components.Velocity.Get(entry)
		var acceleration Vec2.Vec2 = frc.Force.Mult(mss.InverseMass)
		deltaTime := float64(e.Time.DeltaTime().Seconds())
		vel.Velocity.AddUpdate(acceleration.Mult(deltaTime))
		frc.Force.X = 0
		frc.Force.Y = 0
	}
}
