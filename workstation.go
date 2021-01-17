package workstation

import (
	"context"
	"sync"
	"time"
)

func BuildWorkstation(context context.Context, worker Workerable) Stationable {
	return &Workstation{
		mu:        sync.RWMutex{},
		processes: map[string]chan struct{}{},
		worker:    worker,
		context:   context,
		wg:        sync.WaitGroup{},
		done:      make(chan struct{}),
	}
}

// implement Cancellable interface
func (self *Workstation) GetIsCancelledChannel(key string) <-chan struct{} {
	self.mu.RLock()
	ch := self.processes[key]
	self.mu.RUnlock()

	return ch
}

// implement Providable interface
func (self *Workstation) ProvideExecutionContext() context.Context {
	return self.context
}

// implement Stationable interface
func (self *Workstation) PerformAsync(key string, payload Payload) error {
	if self.LookupProcess(key) {
		return ErrorAsyncProcessAlreadyExists
	}

	self.attach(key)

	go func(process string) {
		self.wg.Add(1)

		defer func() {
			// try detach
			self.close(key)
			self.detach(key)
			self.wg.Done()
		}()

		// The method must work synchronously, otherwise it will be completed
		self.worker.Perform(self, process, payload)
	}(key)

	return nil
}

func (self *Workstation) RevokeAsync(key string) error {
	if !self.LookupProcess(key) {
		return ErrorAsyncProcessNotFoundOrAlreadyCompleted
	}

	self.cancel(key)

	return nil
}

// Completion occurs synchronously,
// which represents the possibility of waiting for the completion of all asynchronous processes,
// or an emergency termination
func (self *Workstation) Shutdown() error {
	err := ErrorShutdownWithoutGracefulCompletion

	go func() {
		// wait all async processes
		self.wg.Wait()
		// to inform about the successful completion of the task
		self.done <- struct{}{}
		// without error
		err = nil
	}()

	// waiting for a message about the completion of processes, or completing
	self.await(time.Second * 5)

	return err
}

func (self *Workstation) CountAsync() int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return len(self.processes)
}

// private internal workstation API

func (self *Workstation) LookupProcess(key string) bool {
	self.mu.RLock()
	_, ok := self.processes[key]
	self.mu.RUnlock()

	return ok
}

func (self *Workstation) attach(key string) {
	self.mu.Lock()
	self.processes[key] = make(chan struct{}, 1)
	self.mu.Unlock()
}

// Safe deletion from the pool
func (self *Workstation) detach(key string) {
	self.mu.Lock()
	delete(self.processes, key)
	self.mu.Unlock()
}

func (self *Workstation) cancel(key string) {
	self.mu.Lock()
	defer self.mu.Unlock()

	// Task is cancelled
	self.processes[key] <- struct{}{}
}

func (self *Workstation) close(key string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	close(self.processes[key])
}

func (self *Workstation) channel(key string) <-chan struct{} {
	self.mu.Lock()
	defer self.mu.Unlock()

	return self.processes[key]
}

// The method waits for graceful completion or crashes after a certain amount of time
func (self *Workstation) await(waitDuration time.Duration) {
	select {
	case <-self.done:
		// true
	case <-time.After(waitDuration):
		// false
	}
}
