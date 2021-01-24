package workstation

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func BuildWorkstation(ctx context.Context, workers ...Workerable) *Workstation {
	w := &Workstation{}
	w.mu = sync.RWMutex{}
	w.startedAt = time.Now()
	w.spaces = map[string]*Workspace{}

	for _, worker := range workers {
		w.spaces[worker.Name()] = CreateWorkspace(ctx, worker)
	}

	return w
}

func (w *Workstation) Workspace(name string) (*Workspace, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	workspace, ok := w.spaces[name]

	if !ok {
		return nil, fmt.Errorf("Workspace '%s' not found", name)
	}

	return workspace, nil
}

func (w *Workstation) Shutdown(onError func(err error)) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, space := range w.spaces {
		if err := space.Shutdown(); err != nil {
			onError(err)
		}
	}
}
