package main

import (
	"physengine/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	scene scenes.MyScene
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.scene.Draw(screen)
}

func (g *Game) Update() error {
	g.scene.Update()
	return nil
}

func (g *Game) Layout(width, height int) (int, int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("ECS game")
	if err := ebiten.RunGame(&Game{scenes.MyScene{}}); err != nil {
		panic(err)
	}
}
