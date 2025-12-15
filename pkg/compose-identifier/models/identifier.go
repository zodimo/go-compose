package models

import (
	"github.com/zodimo/go-compose/pkg/compose-identifier/key"
)

type Identifier struct {
	key key.Key
}

func (i Identifier) Value() uint32 {
	return i.key.Value()
}
func (i Identifier) String() string {
	return i.key.String()
}

func NewIdentifier() *Identifier {
	keyManager := key.GetOrCreateKeyManager("default")
	return &Identifier{
		key: keyManager.GenerateKey(),
	}
}

func NewScopedIdentifier(scope string) *Identifier {
	keyManager := key.GetOrCreateKeyManager(scope)
	return &Identifier{
		key: keyManager.GenerateKey(),
	}
}
func NewKeyedIdentifier(scope string, seed string) *Identifier {
	keyManager := key.GetOrCreateKeyManager(scope)
	return &Identifier{
		key: keyManager.CreateKey(seed),
	}
}

type IdentityManager struct {
	scope string
}

func GetOrCreateIdentityManager(scope string) *IdentityManager {
	return &IdentityManager{
		scope: scope,
	}
}

func (im IdentityManager) GenerateKey() *Identifier {
	return NewScopedIdentifier(im.scope)
}
func (im IdentityManager) Create(seed string) *Identifier {
	return NewKeyedIdentifier(im.scope, seed)
}
func (im IdentityManager) ResetKeyCounter() {
	im.getKeyManager().ResetKeyCounter()
}
func (im IdentityManager) EmptyIdentifier() *Identifier {
	return &Identifier{key: im.getKeyManager().EmptyKey()}
}

func (im IdentityManager) getKeyManager() key.KeyManager {
	return key.GetOrCreateKeyManager(im.scope)
}
