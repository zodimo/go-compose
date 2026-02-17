package zipper

import "github.com/zodimo/go-compose/pkg/api"

func NewComposer(options ...api.ComposerOption) Composer {
	opts := api.DefaultComposerOptions()
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	idManager := GetScopedIdentityManager("composer")

	c := &composer{
		timeNowFunc:    opts.TimeNowFunc,
		focus:          nil,
		path:           []pathItem{},
		state:          opts.Store,
		idManager:      idManager,
		locals:         make(map[interface{}]interface{}),
		providersStack: []map[interface{}]interface{}{},
	}

	//wire up the frame lifecycle triggers
	opts.StartFrameTriggerReceiver(c.StartFrame)
	opts.EndFrameTriggerReceiver(c.EndFrame)

	return c
}
