package components

import (
	Vec2 "physengine/helpers/vec2"

	"github.com/yohamta/donburi"
)

type PhysicsBodyData struct {
	Velocity Vec2.Vec2
	AngularVelocity float64
	Force Vec2.Vec2
	Torque float64
	Mass float64
	InvertedMass float64
}

var PhysicsBody = donburi.NewComponentType[PhysicsBodyData]()

func ApplyImpulse(e *donburi.Entry, impulse Vec2.Vec2){
	physbody := PhysicsBody.Get(e)
	change_vel := impulse.Mult(physbody.InvertedMass)
	physbody.Velocity.AddUpdate(change_vel)
}

func ApplyTorque(e *donburi.Entry, torque float64){
	physbody := PhysicsBody.Get(e)
	physbody.Torque += torque
}

func ApplyForce(e *donburi.Entry, force Vec2.Vec2){
	physbody := PhysicsBody.Get(e)
	physbody.Force.AddUpdate(force)
}