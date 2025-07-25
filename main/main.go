package main

import (
	"physengine/components"
	"physengine/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &game.Game{Screen_width: 640, Screen_height: 480, Current_scene: &game.Scene{}}
	ebiten.SetWindowTitle("Game")
	ebiten.SetWindowSize(640, 480)
	g.Current_scene.Root = &game.GameObject{}
	g.Current_scene.Root.Components = append(g.Current_scene.Root.Components, &components.TransformComponent{})
	img := &game.GameObject{}
	img_comp := &components.ImageComponent{}
	img_comp.SetImage("example.png")
	img.Components = append(img.Components, &components.TransformComponent{}, img_comp)
	g.Current_scene.Root.Children = append(g.Current_scene.Root.Children, img)
	g.Current_scene.Root.Init()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
