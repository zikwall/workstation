package workstation

import (
	"context"
)

type (
	Workerable interface {
		Perform(instance Instantiable, key string, payload Payload)
	}
	Observable interface {
		LookupProcess(key string) bool
	}
	Providable interface {
		ProvideExecutionContext() context.Context
	}
	// By implementing this interface you can track the status of manual deletion from the task pool
	// isCanceled := instance.IsCanceled(id)
	// select {
	//		case <-isCanceled:
	// ...
	Cancelable interface {
		GetIsCancelledChannel(key string) <-chan struct{}
	}
	// In order for asynchronous tasks to have access to state and context, the appropriate interfaces are defined
	Instantiable interface {
		Observable
		Providable
		Cancelable
	}
	Stationable interface {
		Instantiable
		PerformAsync(key string, payload Payload) error
		RevokeAsync(key string) error
		CountAsync() int
		Shutdown() error
	}
)
