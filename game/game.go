package game

import (
	"physengine/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Screen_width  int
	Screen_height int
	Current_scene *scene.Scene
	Current_camera
	world *ebiten.Image
}

func (g *Game) Update() error {
	g.Current_scene.Root.ManageUpdate(g)
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Current_scene.Root.ManageDraw(g.world)
	screen.DrawImage(g.world)
}
func (g *Game) Layout(outisdeWidth, outisdeHeight int) (int, int) {
	return g.Screen_width, g.Screen_height
}
