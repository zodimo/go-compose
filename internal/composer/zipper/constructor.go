package zipper

func NewComposer(state PersistentState, options ...ComposerOption) Composer {
	opts := DefaultComposerOptions()
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	idManager := GetScopedIdentityManager("composer")
	idManager.ResetKeyCounter()

	c := &composer{
		timeNowFunc:    opts.TimeNowFunc,
		focus:          nil,
		path:           []pathItem{},
		state:          state,
		idManager:      idManager,
		locals:         make(map[interface{}]interface{}),
		providersStack: []map[interface{}]interface{}{},
	}

	return c
}
