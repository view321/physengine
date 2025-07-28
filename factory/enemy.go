package factory

import (
	"math/rand"
	"physengine/assets"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateEnemy(ecs *ecs.ECS) {
	ecs.World.CreateMany(5, components.Transform, components.BulletTag, components.Velocity, components.CircleCollider, components.Drawable)

	for entry := range components.BulletTag.Iter(ecs.World) {
		donburi.SetValue(entry, components.Velocity, components.VelocityData{X: (rand.Float64() - 0.5) * 10, Y: (rand.Float64() - 0.5) * 10})
		tr := donburi.Get[components.TransformData](entry, components.Transform)
		velocity := components.Velocity.Get(entry)
		velocity.X = (rand.Float64() - 0.5) * 100
		velocity.Y = (rand.Float64() - 0.5) * 100
		tr.Pos.X = (rand.Float64() - 0.5) * 10
		tr.Pos.Y = (rand.Float64() - 0.5) * 10
		tr.Scale = Vec2.Vec2{X: 1, Y: 1}
		dr := components.Drawable.Get(entry)
		img, _ := assets.GetImage("D:/Coding/physengine/assets/assets/enemy.png")
		dr.Sprite = img
		cc := components.CircleCollider.Get(entry)
		cc.Radius = 100
	}
}
