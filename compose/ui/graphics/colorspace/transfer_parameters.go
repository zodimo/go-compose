package colorspace

import (
	"math"
)

// TransferParameters defines the parameters for the ICC parametric curve type 4, as defined in ICC.1:2004-10, section 10.15.
// The corresponding OETF is simply the inverse function.
type TransferParameters struct {
	// Value g in the equation of the EOTF
	Gamma float64
	// Value a in the equation of the EOTF
	A float64
	// Value b in the equation of the EOTF
	B float64
	// Value c in the equation of the EOTF
	C float64
	// Value d in the equation of the EOTF
	D float64
	// Value e in the equation of the EOTF
	E float64
	// Value f in the equation of the EOTF
	F float64
}

const (
	TypePQish  = -2.0
	TypeHLGish = -3.0
)

// NewTransferParameters creates a new TransferParameters instance.
// e and f default to 0.0 if not provided in the Kotlin version, but here we require checking them or use a constructor.
// Since Go doesn't have default arguments, we'll provider helper or just struct literal if safe,
// but validation is complex.
func NewTransferParameters(gamma, a, b, c, d, e, f float64) (TransferParameters, error) {
	if math.IsNaN(a) || math.IsNaN(b) || math.IsNaN(c) || math.IsNaN(d) || math.IsNaN(e) || math.IsNaN(f) || math.IsNaN(gamma) {
		return TransferParameters{}, jsIllegalArgumentException("Parameters cannot be NaN") // checking NaN
	}

	if !isSpecialG(gamma) {
		if !(d >= 0.0 && d <= 1.0) {
			return TransferParameters{}, jsIllegalArgumentException("Parameter d must be in the range [0..1]")
		}

		if d == 0.0 && (a == 0.0 || gamma == 0.0) {
			return TransferParameters{}, jsIllegalArgumentException("Parameter a or g is zero, the transfer function is constant")
		}

		if d >= 1.0 && c == 0.0 {
			return TransferParameters{}, jsIllegalArgumentException("Parameter c is zero, the transfer function is constant")
		}

		if (a == 0.0 || gamma == 0.0) && c == 0.0 {
			return TransferParameters{}, jsIllegalArgumentException("Parameter a or g is zero, and c is zero, the transfer function is constant")
		}

		if c < 0.0 {
			return TransferParameters{}, jsIllegalArgumentException("The transfer function must be increasing")
		}

		if a < 0.0 || gamma < 0.0 {
			return TransferParameters{}, jsIllegalArgumentException("The transfer function must be positive or increasing")
		}
	}

	return TransferParameters{
		Gamma: gamma,
		A:     a,
		B:     b,
		C:     c,
		D:     d,
		E:     e,
		F:     f,
	}, nil
}

func isSpecialG(gamma float64) bool {
	return gamma == TypePQish || gamma == TypeHLGish
}

// Simple error type to mimic IllegalArgumentException
type illegalArgumentError struct {
	s string
}

func (e illegalArgumentError) Error() string {
	return e.s
}

func jsIllegalArgumentException(s string) error {
	return illegalArgumentError{s}
}
