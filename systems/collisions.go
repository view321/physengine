package systems

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateCollisions(e *ecs.ECS){
	resolver_entry, _ := components.CollisionResolverComponent.First(e.World)
	resolver_comp := components.CollisionResolverComponent.Get(resolver_entry)
	for num1 := 0; num1 < len(resolver_comp.Physobs); num1++{
		for num2 := num1+1; num2 < len(resolver_comp.Physobs); num1++{
			ResolveCollisions(resolver_comp.Physobs[num1], resolver_comp.Physobs[num2])
		}
	}
}
func ResolveCollisions(c1, c2 *donburi.Entry){
	
}