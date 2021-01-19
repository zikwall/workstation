package workstation

import (
	"context"
	"sync"
)

type (
	Payload map[string]interface{}
	Process struct {
		ctx    context.Context
		cancel context.CancelFunc
	}
	Workstation struct {
		mu        sync.RWMutex
		processes map[string]Process
		worker    Workerable
		context   context.Context
		// This property simultaneously serves as a counter for asynchronous tasks
		// and a mechanism for waiting/completing the task, for successful completion
		wg sync.WaitGroup
		// This property serves as a flag for successful completion of all asynchronous tasks
		done chan struct{}
	}
)
