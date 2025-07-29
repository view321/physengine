package systems

import (
	"fmt"
	"math"
	"math/rand"
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi/ecs"
)

func UpdateEnemy(ecs *ecs.ECS) {
	for entry := range components.BulletTag.Iter(ecs.World) {
		tr := components.Transform.Get(entry)
		v := components.Velocity.Get(entry)

		if tr.Pos.X > 1000 {
			tr.Pos.X = 1000
			v.Velocity.X *= -1
			v.Velocity.X = v.Velocity.X + (rand.Float64()-0.5)*10
			v.Velocity.Y = v.Velocity.Y + (rand.Float64()-0.5)*10
			magn := math.Sqrt(math.Pow(v.Velocity.X, 2) + math.Pow(v.Velocity.Y, 2))
			v.Velocity.X = v.Velocity.X / magn * 25
			v.Velocity.Y = v.Velocity.Y / magn * 25
		}
		if tr.Pos.X < -1000 {
			tr.Pos.X = -1000
			v.Velocity.X *= -1
			v.Velocity.X = v.Velocity.X + (rand.Float64()-0.5)*10
			v.Velocity.Y = v.Velocity.Y + (rand.Float64()-0.5)*10
			magn := math.Sqrt(math.Pow(v.Velocity.X, 2) + math.Pow(v.Velocity.Y, 2))
			v.Velocity.X = v.Velocity.X / magn * 25
			v.Velocity.Y = v.Velocity.Y / magn * 25
		}
		if tr.Pos.Y > 1000 {
			tr.Pos.Y = 1000
			v.Velocity.Y *= -1
			v.Velocity.X = v.Velocity.X + (rand.Float64()-0.5)*10
			v.Velocity.Y = v.Velocity.Y + (rand.Float64()-0.5)*10
			magn := math.Sqrt(math.Pow(v.Velocity.X, 2) + math.Pow(v.Velocity.Y, 2))
			v.Velocity.X = v.Velocity.X / magn * 25
			v.Velocity.Y = v.Velocity.Y / magn * 25
		}
		if tr.Pos.Y < -1000 {
			tr.Pos.Y = -1000
			v.Velocity.Y *= -1
			v.Velocity.X = v.Velocity.X + (rand.Float64()-0.5)*10
			v.Velocity.Y = v.Velocity.Y + (rand.Float64()-0.5)*10
			magn := math.Sqrt(math.Pow(v.Velocity.X, 2) + math.Pow(v.Velocity.Y, 2))
			v.Velocity.X = v.Velocity.X / magn * 25
			v.Velocity.Y = v.Velocity.Y / magn * 25
		}

		player, _ := components.Player.First(ecs.World)

		if components.CirclesCollide(player, entry) {
			fmt.Println("Game Over")
		}

		components.SetPos(entry, tr.Pos.Add(Vec2.Vec2{X: v.Velocity.X, Y: v.Velocity.Y}))
	}
}
