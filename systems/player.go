package systems

import (
	"physengine/components"
	Vec2 "physengine/helpers/vec2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi/ecs"
)

func UpdatePlayer(ecs *ecs.ECS) {
	playerEntry, _ := components.Player.First(ecs.World)
	tr := components.Transform.Get(playerEntry)
	velocity := components.Velocity.Get(playerEntry)
	velocity.X = 0
	velocity.Y = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		velocity.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		velocity.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		velocity.X += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		velocity.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		components.SetRot(playerEntry, tr.Rot+0.1)
	}
	components.SetPos(playerEntry, Vec2.Vec2{X: tr.Pos.X + velocity.X*ecs.Time.DeltaTime().Seconds()*100, Y: tr.Pos.Y + velocity.Y*ecs.Time.DeltaTime().Seconds()*100})
}
