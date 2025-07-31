package systems

import (
	"physengine/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateCollisions(e *ecs.ECS) {
	settings_entry, _ := components.CollisionSettings.First(e.World)
	settings := components.CollisionSettings.Get(settings_entry)
	for _, manifold := range settings.Manifolds {
		if manifold.Layer1 == manifold.Layer2 {
			for i := 0; i < len(manifold.Layer1.Entries); i++ {
				for j := i + 1; j < len(manifold.Layer2.Entries); j++ {
					TryResolvingCollision(manifold.Layer1.Entries[i], manifold.Layer2.Entries[j], manifold.Settings)
				}
			}
		} else {
			for _, e1 := range manifold.Layer1.Entries {
				for _, e2 := range manifold.Layer2.Entries {
					TryResolvingCollision(e1, e2, manifold.Settings)
				}
			}
		}
	}
}
func TryResolvingCollision(e1, e2 *donburi.Entry, settings components.CollisionSettingsData) {
	if e1.HasComponent(components.PolygonCollider) && e2.HasComponent(components.PolygonCollider) {
		plgn1 := components.PolygonCollider.Get(e1)
		plgn2 := components.PolygonCollider.Get(e2)
	}
}
