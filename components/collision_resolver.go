package components

import (
	"github.com/yohamta/donburi"
)

type CollisionResolverData struct {
	Physobs []*donburi.Entry
}

var CollisionResolverComponent = donburi.NewComponentType[CollisionResolverData]()
