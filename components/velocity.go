package components

import (
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type VelocityData struct{
	Velocity Vec2.Vec2
}

var Velocity = donburi.NewComponentType[VelocityData]()