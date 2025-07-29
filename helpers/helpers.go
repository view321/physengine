package helpers

func Clamp(val, min, max float64) float64{
	if val > max{
		val = max
	}
	if val < min{
		val = min
	}
	return val
}