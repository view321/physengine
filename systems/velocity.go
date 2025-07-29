package systems

import (
	"fmt"
	"physengine/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateVelocity(e *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Velocity))
	for entry := range query.Iter(e.World){
		tr := components.Transform.Get(entry)
		vel := components.Velocity.Get(entry)
		fmt.Println(tr.Pos.X, tr.Pos.Y)
		components.SetPos(entry, tr.Pos.Add(vel.Velocity.Mult(float64(e.Time.DeltaTime().Seconds()))))
	}
}