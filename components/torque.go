package components

import "github.com/yohamta/donburi"

type TorqueData struct {
	Torque float64
}

var Torque = donburi.NewComponentType[TorqueData]()

// SetTorque sets the torque applied to an entity
func SetTorque(entry *donburi.Entry, torque float64) {
	torqueData := Torque.Get(entry)
	if torqueData != nil {
		torqueData.Torque = torque
	}
}

// AddTorque adds to the current torque
func AddTorque(entry *donburi.Entry, torque float64) {
	torqueData := Torque.Get(entry)
	if torqueData != nil {
		torqueData.Torque += torque
	}
}

// GetTorque returns the current torque
func GetTorque(entry *donburi.Entry) float64 {
	torqueData := Torque.Get(entry)
	if torqueData != nil {
		return torqueData.Torque
	}
	return 0.0
} 