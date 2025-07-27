package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type DrawableData struct{
	Sprite *ebiten.Image
}

var Drawable = donburi.NewComponentType[DrawableData]()