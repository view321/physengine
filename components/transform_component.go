package components

import (
	"physengine/game"
	"physengine/vector2"
	"physengine/gameobjects"
	"github.com/hajimehoshi/ebiten/v2"
)

type TransformComponent struct {
	Position   vector2.Vector2
	Rotation   float64
	Scale      vector2.Vector2
	gameObject *gameobjects.GameObject
}

func (*TransformComponent) Update(g *game.Game) error {
	return nil
}
func (*TransformComponent) Draw(screen *ebiten.Image) {

}
func (*TransformComponent) GetName() string {
	return "TransformComponent"
}
func (t *TransformComponent) SetPos(v vector2.Vector2) {
	var diff vector2.Vector2 = vector2.Add(v, t.Position.Mult(-1))
	t.Position = v
	for _, obj := range t.gameObject.Children {
		(*obj.Transform).(*TransformComponent).ChangePos(diff)
	}

}
func (t *TransformComponent) ChangePos(v vector2.Vector2) {
	for _, obj := range t.gameObject.Children {
		(*obj.Transform).(*TransformComponent).ChangePos(v)
	}
}
func (t *TransformComponent) SetScale(n vector2.Vector2) {
	t.Scale = n
	for _, obj := range t.gameObject.Children {
		(*obj.Transform).(*TransformComponent).SetScale(n)
	}
}
func (t *TransformComponent) ScaleBy(n vector2.Vector2) {
	t.Scale = vector2.Vector2{X: t.Scale.X * n.X, Y: t.Scale.Y * n.Y}
	for _, obj := range t.gameObject.Children {
		(*obj.Transform).(*TransformComponent).ScaleBy(n)
	}
}
func (t *TransformComponent) SetRotation(n float64) {
	t.Rotation = n
	for _, obj := range t.gameObject.Children {
		(*obj.Transform).(*TransformComponent).SetRotation(n)
	}
}
func (t *TransformComponent) RotateBy(n float64) {
	t.Rotation = t.Rotation + n
	for _, obj := range t.gameObject.Children {
		(*obj.Transform).(*TransformComponent).RotateBy(n)
	}
}
func (t *TransformComponent) Init(g *gameobjects.GameObject) {
	t.gameObject = g
	if g.Parent == nil {
		t.Position = vector2.Vector2{X: 0, Y: 0}
		t.Rotation = 0
		t.Scale = vector2.Vector2{X: 1, Y: 1}
		return
	}
	t.Position = vector2.Vector2{X: (*g.Parent.Transform).(*TransformComponent).Position.X, Y: (*g.Parent.Transform).(*TransformComponent).Position.Y}
	t.Rotation = (*g.Parent.Transform).(*TransformComponent).Rotation
	t.Scale = (*g.Parent.Transform).(*TransformComponent).Scale
}
