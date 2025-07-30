package components

import (
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type ForceData struct {
	Force Vec2.Vec2
}

var ForceComponent = donburi.NewComponentType[ForceData]()

func SetForce(entry *donburi.Entry, newForce Vec2.Vec2){
	var frc *ForceData = ForceComponent.Get(entry)
	frc.Force = newForce
}

func AddForce(entry *donburi.Entry, force_diff Vec2.Vec2){
	var frc *ForceData = ForceComponent.Get(entry)
	frc.Force.AddUpdate(force_diff)
}