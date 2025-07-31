package components

import (
	"github.com/yohamta/donburi"
)

type CollisionLayer struct {
	Entries []*donburi.Entry
}

type CollisionManifold struct{
	Layer1 *CollisionLayer
	Layer2 *CollisionLayer
	Settings CollisionSettingsData
}

type CollisionSettingsData struct{
	Collide bool
}

type CollisionData struct {
	Manifolds []*CollisionManifold
}

var CollisionSettings = donburi.NewComponentType[CollisionData]()