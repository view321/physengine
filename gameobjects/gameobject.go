package gameobjects

import (
	"errors"
	"physengine/components"
	"physengine/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameObject struct {
	Transform  *components.Component
	Components []components.Component
	Children   []*GameObject
	Parent     *GameObject
}

func (gameobject *GameObject) ManageUpdate(g *game.Game) {
	for _, component := range (*gameobject).Components {
		component.Update(g)
	}
	for _, child := range (*gameobject).Children {
		child.ManageUpdate(g)
	}
}
func (gameobject *GameObject) ManageDraw(screen *ebiten.Image) {
	for _, component := range (*gameobject).Components {
		component.Draw(screen)
	}
	for _, child := range (*gameobject).Children {
		child.ManageDraw(screen)
	}
}
func (gameobject *GameObject) Init() {
	for _, comp := range gameobject.Components {
		comp.Init(gameobject)
	}
	for _, child := range gameobject.Children {
		child.Init()
	}
}
func (gameObject *GameObject) GetComponent(name string) (*components.Component, error) {
	for _, comp := range gameObject.Components {
		if comp.GetName() == name {
			return &comp, nil
		}
	}
	return nil, errors.New("did not find the component")
}
