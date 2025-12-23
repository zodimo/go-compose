package colorspace

import (
	"math"
)

// RcpResponse returns the reciprocal parametric response (OETF key).
// if x >= d*c: (x^(1/g) - b) / a
// else: x / c
func RcpResponse(x, a, b, c, d, g float64) float64 {
	if x >= d*c {
		return (math.Pow(x, 1.0/g) - b) / a
	}
	return x / c
}

// Response returns the parametric response (EOTF key).
// if x >= d: (a*x + b)^g
// else: c*x
func Response(x, a, b, c, d, g float64) float64 {
	if x >= d {
		return math.Pow(a*x+b, g)
	}
	return c * x
}

// RcpResponse extended with e and f
func RcpResponseExtended(x, a, b, c, d, e, f, g float64) float64 {
	if x >= d*c {
		return (math.Pow(x-e, 1.0/g) - b) / a
	}
	return (x - f) / c
}

// Response extended with e and f
func ResponseExtended(x, a, b, c, d, e, f, g float64) float64 {
	if x >= d {
		return math.Pow(a*x+b, g) + e
	}
	return c*x + f
}

// AbsRcpResponse handles negative values by processing abs(x) and restoring sign
func AbsRcpResponse(x, a, b, c, d, g float64) float64 {
	sign := 1.0
	if x < 0 {
		sign = -1.0
	}
	return RcpResponse(math.Abs(x), a, b, c, d, g) * sign
}

// AbsResponse handles negative values by processing abs(x) and restoring sign
func AbsResponse(x, a, b, c, d, g float64) float64 {
	sign := 1.0
	if x < 0 {
		sign = -1.0
	}
	return Response(math.Abs(x), a, b, c, d, g) * sign
}
