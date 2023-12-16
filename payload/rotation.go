package payload

import "math"

func convertRotationToMatrix(degrees float64) []float64 {
	// ! first convert degrees to radians
	radians := degrees * math.Pi / 180
	// ! now we do this:
	/*
		[cos(r), -sin(r), 0]
		[sin(r), cos(r), 0]

		in this order:
		[cos(r), sin(r), -sin(r), cos(r), 0, 0]
	*/
	c := math.Cos(radians)
	s := math.Sin(radians)
	return []float64{c, s, -s, c, 0, 0}
}
