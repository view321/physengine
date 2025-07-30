package systems

import (
	"physengine/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateTorque(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Torque, components.AngularVelocity, components.MassComponent))
	for entry := range query.Iter(e.World) {
		torque := components.GetTorque(entry)
		mass := components.MassComponent.Get(entry)
		
		if mass != nil && mass.InverseInertia > 0 {
			// Calculate angular acceleration: α = τ / I
			angularAcceleration := torque * mass.InverseInertia
			
			// Update angular velocity: ω = ω₀ + α * dt
			deltaTime := float64(e.Time.DeltaTime().Seconds())
			angularVelocityDelta := angularAcceleration * deltaTime
			
			components.ChangeAngularVelocity(entry, angularVelocityDelta)
		}
		
		// Reset torque for next frame
		components.SetTorque(entry, 0.0)
	}
} 