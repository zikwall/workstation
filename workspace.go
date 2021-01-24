package workstation

import (
	"context"
	"sync"
	"time"
)

func CreateWorkspace(context context.Context, worker Workerable) *Workspace {
	return &Workspace{
		mu:        sync.RWMutex{},
		processes: map[string]Process{},
		worker:    worker,
		context:   context,
		wg:        sync.WaitGroup{},
		done:      make(chan struct{}),
	}
}

// implement Stationable interface
func (self *Workspace) PerformAsync(key string, payload Payload) error {
	if self.LookupProcess(key) {
		return ErrorAsyncProcessAlreadyExists
	}

	ctx, cancel := context.WithCancel(self.context)

	self.attach(key, Process{
		ctx:    ctx,
		cancel: cancel,
	})

	go func(process string) {
		self.wg.Add(1)

		defer func() {
			self.tryCancelAndDetach(key)
			self.wg.Done()
		}()

		// The method must work synchronously, otherwise it will be completed
		self.worker.Perform(ctx, process, payload)
	}(key)

	return nil
}

func (self *Workspace) RevokeAsync(key string) error {
	if !self.LookupProcess(key) {
		return ErrorAsyncProcessNotFoundOrAlreadyCompleted
	}

	self.cancel(key)
	self.detach(key)

	return nil
}

func (self *Workspace) CountAsync() int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return len(self.processes)
}

func (self *Workspace) LookupProcess(key string) bool {
	self.mu.RLock()
	_, ok := self.processes[key]
	self.mu.RUnlock()

	return ok
}

// Completion occurs synchronously,
// which represents the possibility of waiting for the completion of all asynchronous processes,
// or an emergency termination
func (self *Workspace) Shutdown() error {
	err := ErrorShutdownWithoutGracefulCompletion

	go func() {
		// wait all async processes
		self.wg.Wait()
		// without error
		err = nil
		// to inform about the successful completion of the task
		self.done <- struct{}{}
	}()

	// waiting for a message about the completion of processes, or completing
	self.await(time.Second * 5)

	return err
}

// private internal workstation API

func (self *Workspace) attach(key string, process Process) {
	self.mu.Lock()
	self.processes[key] = process
	self.mu.Unlock()
}

// Safe deletion from the pool
func (self *Workspace) detach(key string) {
	self.mu.Lock()
	delete(self.processes, key)
	self.mu.Unlock()
}

func (self *Workspace) cancel(key string) {
	self.mu.Lock()
	defer self.mu.Unlock()

	// Task is cancelled
	self.processes[key].cancel()
}

func (self *Workspace) tryCancelAndDetach(key string) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if _, ok := self.processes[key]; ok && self.processes[key].ctx.Err() == nil {
		self.processes[key].cancel()
	}

	delete(self.processes, key)
}

// The method waits for graceful completion or crashes after a certain amount of time
func (self *Workspace) await(waitDuration time.Duration) {
	select {
	case <-self.done:
		// true
	case <-time.After(waitDuration):
		// false
	}
}
