package payload

func interpolate(from, to []float64, f float64) []float64 {
	out := []float64{}
	for i := 0; i < len(from); i++ {
		out = append(out, interpolateNum(from[i], to[i], f))
	}
	return out
}

func interpolateNum(from, to, f float64) float64 {
	return from*(1-f) + to*f
}
