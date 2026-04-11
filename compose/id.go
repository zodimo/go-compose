package compose

// Identity Composable
func Id() Composable {
	return func(c Composer) Composer {
		return c
	}
}
