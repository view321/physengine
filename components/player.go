package components

import (
	"github.com/yohamta/donburi"
)

type PlayerData struct{
	Lives int64
}

var Player = donburi.NewComponentType[PlayerData]()