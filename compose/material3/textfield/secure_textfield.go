package textfield

const (
	DefaultObfuscationCharacter = '•'
	SecureTextFieldNodeID       = "Material3SecureTextField"
)

// SecureTextField implements the Filled Secure Text Field.
// It is designed for password entry.
func SecureTextField(
	value string,
	onValueChange func(string),
	options ...TextFieldOption,
) Composable {
	defaultOpts := []TextFieldOption{
		WithMask(DefaultObfuscationCharacter),
		WithSingleLine(true),
	}
	mergedOptions := append(defaultOpts, options...)
	return Filled(value, onValueChange, mergedOptions...)
}

// OutlinedSecureTextField implements the Outlined Secure Text Field.
func OutlinedSecureTextField(
	value string,
	onValueChange func(string),
	options ...TextFieldOption,
) Composable {
	defaultOpts := []TextFieldOption{
		WithMask(DefaultObfuscationCharacter),
		WithSingleLine(true),
	}
	mergedOptions := append(defaultOpts, options...)
	return Outlined(value, onValueChange, mergedOptions...)
}
