package colorspace

// Adaptation represents a chromatic adaptation matrix used for the von Kries transform.
// These matrices are used to convert values in the CIE XYZ space to values in the LMS space.
type Adaptation struct {
	Transform []float32
}

var (
	// AdaptationBradford is the Bradford chromatic adaptation transform, as defined in the CIECAM97s color appearance model.
	AdaptationBradford = Adaptation{
		Transform: []float32{
			0.8951, -0.7502, 0.0389,
			0.2664, 1.7135, -0.0685,
			-0.1614, 0.0367, 1.0296,
		},
	}

	// AdaptationVonKries is the von Kries chromatic adaptation transform.
	AdaptationVonKries = Adaptation{
		Transform: []float32{
			0.40024, -0.22630, 0.00000,
			0.70760, 1.16532, 0.00000,
			-0.08081, 0.04570, 0.91822,
		},
	}

	// AdaptationCiecat02 is the CIECAT02 chromatic adaption transform, as defined in the CIECAM02 color appearance model.
	AdaptationCiecat02 = Adaptation{
		Transform: []float32{
			0.7328, -0.7036, 0.0030,
			0.4296, 1.6975, 0.0136,
			-0.1624, 0.0061, 0.9834,
		},
	}
)

func (a Adaptation) String() string {
	switch {
	// Compare slices by value or identity? Identity for these globals is fine.
	// But in Go slice comparison is not directly possible with ==.
	// We'll rely on simple matching or just not implementing String() based on value deeply.
	// For now, let's just return "Adaptation" or check pointers if we use pointers.
	// Since we use structs with slices, we can't easily compare equality without helper.
	// Let's assume standard ones are used.
	}
	return "Adaptation"
}
