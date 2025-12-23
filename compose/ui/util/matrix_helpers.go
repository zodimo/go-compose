package util

// Mul3x3Float3 multiplies a 3x3 matrix (row-major flat array) by a 3-component vector.
// Result maps back to v.
func Mul3x3Float3(m []float32, v []float32) []float32 {
	x, y, z := v[0], v[1], v[2]
	v[0] = m[0]*x + m[3]*y + m[6]*z
	v[1] = m[1]*x + m[4]*y + m[7]*z
	v[2] = m[2]*x + m[5]*y + m[8]*z
	return v
}

// Inverse3x3 returns the inverse of a 3x3 matrix.
func Inverse3x3(m []float32) []float32 {
	a, b, c := m[0], m[3], m[6]
	d, e, f := m[1], m[4], m[7]
	g, h, i := m[2], m[5], m[8]

	A := e*i - f*h
	B := f*g - d*i
	C := d*h - e*g

	det := a*A + b*B + c*C
	// If det is 0, we can't invert. Return zero matrix or panic?
	// Kotlin version assumes invertible.

	invDet := 1.0 / det

	res := make([]float32, 9)
	res[0] = A * invDet
	res[1] = B * invDet
	res[2] = C * invDet
	res[3] = (c*h - b*i) * invDet
	res[4] = (a*i - c*g) * invDet
	res[5] = (b*g - a*h) * invDet
	res[6] = (b*f - c*e) * invDet
	res[7] = (c*d - a*f) * invDet
	res[8] = (a*e - b*d) * invDet
	return res
}

// Mul3x3 multiplies two 3x3 matrices (lhs * rhs).
func Mul3x3(lhs, rhs []float32) []float32 {
	r := make([]float32, 9)
	r[0] = lhs[0]*rhs[0] + lhs[3]*rhs[1] + lhs[6]*rhs[2]
	r[1] = lhs[1]*rhs[0] + lhs[4]*rhs[1] + lhs[7]*rhs[2]
	r[2] = lhs[2]*rhs[0] + lhs[5]*rhs[1] + lhs[8]*rhs[2]
	r[3] = lhs[0]*rhs[3] + lhs[3]*rhs[4] + lhs[6]*rhs[5]
	r[4] = lhs[1]*rhs[3] + lhs[4]*rhs[4] + lhs[7]*rhs[5]
	r[5] = lhs[2]*rhs[3] + lhs[5]*rhs[4] + lhs[8]*rhs[5]
	r[6] = lhs[0]*rhs[6] + lhs[3]*rhs[7] + lhs[6]*rhs[8]
	r[7] = lhs[1]*rhs[6] + lhs[4]*rhs[7] + lhs[7]*rhs[8]
	r[8] = lhs[2]*rhs[6] + lhs[5]*rhs[7] + lhs[8]*rhs[8]
	return r
}
