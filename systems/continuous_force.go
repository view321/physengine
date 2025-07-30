package systems

import (
	"physengine/components"

	"github.com/yohamta/donburi/ecs"
)

func ApplyContinuousForce(e *ecs.ECS) {
	cont_entry, _ := components.ContForceComp.First(e.World)
	cont_force_comp := components.ContForceComp.Get(cont_entry)
	for _, cont := range cont_force_comp.Containers {
		for _, obj := range cont.Objects {
			components.AddForce(obj, cont.Force)
		}
	}
}
