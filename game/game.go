package game

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Screen_width  int
	Screen_height int
	Current_scene *Scene
}

func (g *Game) Update() error {
	g.Current_scene.Root.ManageUpdate(g)
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Current_scene.Root.ManageDraw(screen)
}
func (g *Game) Layout(outisdeWidth, outisdeHeight int) (int, int) {
	return g.Screen_width, g.Screen_height
}

type GameObject struct {
	Transform  *Component
	Components []Component
	Children   []*GameObject
	Parent     *GameObject
}

func (gameobject *GameObject) ManageUpdate(g *Game) {
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
func (gameObject *GameObject) GetComponent(name string) (*Component, error) {
	for _, comp := range gameObject.Components {
		if comp.GetName() == name {
			return &comp, nil
		}
	}
	return nil, errors.New("did not find the component")
}

type Component interface {
	Update(g *Game) error
	Draw(screen *ebiten.Image)
	GetName() string
	Init(g *GameObject)
}
type Scene struct {
	Root *GameObject
}
