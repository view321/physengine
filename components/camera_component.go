package components

import (
	"physengine/gameobjects"
	"physengine/vector2"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
)

type CameraComponent struct {
	Viewport   vector2.Vector2
	Active     bool
	ZoomFactor int
	transform  *TransformComponent
	gameObject *gameobjects.GameObject
}

func (c *CameraComponent) viewportCenter() vector2.Vector2 {
	return vector2.Vector2{X: c.transform.Position.X + c.Viewport.X/2, Y: c.transform.Position.Y + c.Viewport.Y/2}
}
func (c *CameraComponent) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.viewportCenter().X, -c.viewportCenter().Y)
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(c.transform.Rotation * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter().X, c.viewportCenter().Y)
	return m
}
func (c *CameraComponent) Update() error {
	return nil
}
func (c *CameraComponent) Init(g *gameobjects.GameObject) {
	c.gameObject = g
	for _, comp := range g.Components {
		if comp.GetName() == "TransformComponent" {
			c.transform = comp.(*components.TransformComponent)
		}
	}
}
func (c *CameraComponent) Draw(screen *ebiten.Image) {

}
func (c *CameraComponent) GetName() string {
	return "CameraComponent"
}
