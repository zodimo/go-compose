package identity

type IdentityManager interface {
	GenerateID() Identifier
	CreateID(seed string) Identifier
	ResetKeyCounter()
	EmptyIdentifier() Identifier
	private()
}
