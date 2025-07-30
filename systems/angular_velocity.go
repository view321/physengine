package systems

import (
	"physengine/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateAngularVelocity(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.AngularVelocity))
	for entry := range query.Iter(e.World) {
		angVel := components.AngularVelocity.Get(entry)
		
		// Update rotation based on angular velocity and delta time
		rotationDelta := angVel.AngularVelocity * float64(e.Time.DeltaTime().Seconds())
		components.Rotate(entry, rotationDelta)
	}
} 