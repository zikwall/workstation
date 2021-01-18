package workstation

import (
	"context"
)

type (
	Workerable interface {
		Perform(ctx context.Context, key string, payload Payload)
	}
	Observable interface {
		LookupProcess(key string) bool
		CountAsync() int
	}
	// In order for asynchronous tasks to have access to state and context, the appropriate interfaces are defined
	Instantiable interface {
		Observable
	}
	Stationable interface {
		Instantiable
		PerformAsync(key string, payload Payload) error
		RevokeAsync(key string) error
		Shutdown() error
	}
)
