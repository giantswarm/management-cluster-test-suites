package state

import (
	"context"
	"sync"

	"github.com/giantswarm/clustertest/v2"
)

var lock = &sync.Mutex{}

type state struct {
	framework *clustertest.Framework
	ctx       context.Context
}

var singleInstance *state

func get() *state {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &state{}
		}
	}

	return singleInstance
}

func SetContext(ctx context.Context) {
	s := get()
	s.ctx = ctx
}

func GetContext() context.Context {
	return get().ctx
}

func SetFramework(framework *clustertest.Framework) {
	s := get()
	s.framework = framework
}

func GetFramework() *clustertest.Framework {
	return get().framework
}
