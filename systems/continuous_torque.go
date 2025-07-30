package systems

import (
	"physengine/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func ApplyContinuousTorque(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Torque, components.AngularVelocity, components.MassComponent))
	for entry := range query.Iter(e.World) {
		// Get the current torque and add a continuous torque
		currentTorque := components.GetTorque(entry)
		components.SetTorque(entry, currentTorque + 1000.0) // Add continuous torque
	}
} 