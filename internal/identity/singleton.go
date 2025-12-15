package identity

var _ IdentityManager = (*identityManager)(nil)

type identityManager struct {
	manager apiIdentityManager
}

func (im identityManager) GenerateID() Identifier {
	return im.manager.GenerateKey()
}
func (im identityManager) ResetKeyCounter() {
	im.manager.ResetKeyCounter()
}
func (im identityManager) EmptyIdentifier() Identifier {
	return im.manager.EmptyIdentifier()
}

func (im identityManager) CreateID(seed string) Identifier {
	return im.manager.Create(seed)
}
func (im identityManager) private() {}

func newIdentityManager(manager apiIdentityManager) IdentityManager {
	return identityManager{
		manager: manager,
	}
}

func GetIdentityManager() IdentityManager {
	manager := getOrCreateIdentityManager("default")
	return newIdentityManager(manager)
}

func GetScopedIdentityManager(scope string) IdentityManager {
	manager := getOrCreateIdentityManager(scope)
	return newIdentityManager(manager)
}
