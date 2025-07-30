package components

import (
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type ContForceContainer struct {
	Force Vec2.Vec2
	Objects []*donburi.Entry
}

type ContForceData struct{
	Containers []ContForceContainer
}

var ContForceComp = donburi.NewComponentType[ContForceData]()