package components

import (
	"image"
	_ "image/png"
	"os"
	"physengine/gameobjects"
	"physengine/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageComponent struct {
	Img        *ebiten.Image
	transform  *TransformComponent
	gameObject *gameobjects.GameObject
}

func (ic *ImageComponent) Init(g *gameobjects.GameObject) {
	ic.gameObject = g
	for _, comp := range g.Components {
		if comp.GetName() == "TransformComponent" {
			ic.transform = comp.(*TransformComponent)
		}
	}
}
func (ic *ImageComponent) SetImage(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	ic.Img = ebiten.NewImageFromImage(img)
}
func (ic *ImageComponent) Draw(screen *ebiten.Image) {
	var options ebiten.DrawImageOptions
	options.GeoM.Translate(ic.transform.Position.X, -ic.transform.Position.Y)
	options.GeoM.Rotate(ic.transform.Rotation)
	options.GeoM.Scale(ic.transform.Scale.X, ic.transform.Scale.Y)
	screen.DrawImage(ic.Img, &options)
}
func (ic *ImageComponent) Update(g *game.Game) error {
	return nil
}
func (ic *ImageComponent) GetName() string {
	return "ImageComponent"
}
