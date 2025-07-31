package systems

import (
	"physengine/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdatePhysics(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.PhysicsBody))
	for entry := range query.Iter(e.World) {
		physbody := components.PhysicsBody.Get(entry)
		seconds := e.Time.DeltaTime().Seconds()
		components.ChangePos(entry, physbody.Velocity.Mult(seconds))
		physbody.Velocity.AddUpdate(physbody.Force.Mult(physbody.InvertedMass))
		components.Rotate(entry, physbody.AngularVelocity * seconds)
		physbody.AngularVelocity += (physbody.Torque * physbody.InvertedMass)
	}
}
