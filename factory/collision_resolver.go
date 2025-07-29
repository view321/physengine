package factory

import (
	"physengine/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func CreateCollisionResolver(e *ecs.ECS) *donburi.Entry{
	entity := e.World.Create(components.CollisionResolverComponent)
	entry := e.World.Entry(entity)
	resolve_comp := components.CollisionResolverComponent.Get(entry)
	query := donburi.NewQuery(filter.Or(filter.Contains(components.CircleCollider), filter.Contains(components.AABB_Component)))
	for phys_entry := range query.Iter(e.World){
		resolve_comp.Physobs = append(resolve_comp.Physobs, phys_entry)
	}
	return entry
}