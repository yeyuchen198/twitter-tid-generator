package payload

type cubic struct {
	Curves [4]float64
}

func (c *cubic) getValue(time float64) float64 {
	startGradient := 0.0
	endGradient := 0.0
	if time <= 0.0 {
		if c.Curves[0] > 0.0 {
			startGradient = c.Curves[1] / c.Curves[0]
		} else if c.Curves[1] == 0.0 && c.Curves[2] > 0.0 {
			startGradient = c.Curves[3] / c.Curves[2]
		}
		return startGradient * time
	}

	if time >= 1.0 {
		if c.Curves[2] < 1.0 {
			endGradient = (c.Curves[3] - 1.0) / (c.Curves[2] - 1.0)
		} else if c.Curves[2] == 1.0 && c.Curves[0] < 1.0 {
			endGradient = (c.Curves[1] - 1.0) / (c.Curves[0] - 1.0)
		}
		return 1.0 + endGradient*(time-1.0)
	}

	start := 0.0
	end := 1.0
	mid := 0.0
	for start < end {
		mid = (start + end) / 2
		xEst := f(c.Curves[0], c.Curves[2], mid)
		if abs(time-xEst) < 0.00001 {
			return f(c.Curves[1], c.Curves[3], mid)
		}
		if xEst < time {
			start = mid
		} else {
			end = mid
		}
	}
	return f(c.Curves[1], c.Curves[3], mid)
}

func abs(in float64) float64 {
	if in < 0 {
		return -in
	}
	return in
}

func f(a, b, m float64) float64 {
	return 3.0*a*(1-m)*(1-m)*m + 3.0*b*(1-m)*m*m + m*m*m
}
