package components

import "github.com/yohamta/donburi"

type AngularVelocityData struct {
	AngularVelocity float64
}

var AngularVelocity = donburi.NewComponentType[AngularVelocityData]()

// SetAngularVelocity sets the angular velocity of an entity
func SetAngularVelocity(entry *donburi.Entry, angularVelocity float64) {
	angVel := AngularVelocity.Get(entry)
	if angVel != nil {
		angVel.AngularVelocity = angularVelocity
	}
}

// ChangeAngularVelocity adds to the current angular velocity
func ChangeAngularVelocity(entry *donburi.Entry, angularVelocityDelta float64) {
	angVel := AngularVelocity.Get(entry)
	if angVel != nil {
		angVel.AngularVelocity += angularVelocityDelta
	}
}

// GetAngularVelocity returns the current angular velocity
func GetAngularVelocity(entry *donburi.Entry) float64 {
	angVel := AngularVelocity.Get(entry)
	if angVel != nil {
		return angVel.AngularVelocity
	}
	return 0.0
} 