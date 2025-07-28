package systems

import (
	"physengine/components"

	"github.com/yohamta/donburi/ecs"
)

func UpdatePlayer(ecs *ecs.ECS) {
	playerEntry, _ := components.Player.First(ecs.World)
	cam, _ := components.Camera.First(ecs.World)
	cam_comp := components.Camera.Get(cam)
	components.SetPos(playerEntry, cam_comp.LastMousePos)
}
