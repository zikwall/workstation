package workstation

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// sample
type Collector struct {
	mu sync.RWMutex
	c  []interface{}
}

func (c *Collector) Add(s ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c = append(c.c, s...)
}

func (c *Collector) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.c)
}

func (c *Collector) All() []interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.c
}

var globalCollector = &Collector{
	mu: sync.RWMutex{},
	c:  []interface{}{},
}

type MockWorker struct{}

func (w *MockWorker) Perform(instance Instantiable, key string, payload Payload) {
	ctx, cancel := context.WithCancel(instance.ProvideExecutionContext())

	defer cancel()

	isCanceled := instance.GetIsCancelledChannel(key)

	for instance.ObserveProcessAlive(key) {
		select {
		case <-ctx.Done():
			return
		case <-isCanceled:
			return
		default:
			globalCollector.Add(payload["id"])
		}
	}
}

func TestWorkstation(t *testing.T) {
	t.Run("it should be successful init workstation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		workstation := BuildWorkstation(ctx, &MockWorker{})

		t.Run("it should be success create three processes", func(t *testing.T) {
			if err := workstation.PerformAsync("process_one", Payload{"id": 10, "name": "Process One"}); err != nil {
				t.Fatal(err)
			}

			if err := workstation.PerformAsync("process_two", Payload{"id": 20, "name": "Process Two"}); err != nil {
				t.Fatal(err)
			}

			if err := workstation.PerformAsync("process_three", Payload{"id": 30, "name": "Process Three"}); err != nil {
				t.Fatal(err)
			}

			t.Run("it should be success check count of active process", func(t *testing.T) {
				if workstation.CountAsync() != 3 {
					t.Fatal("Failed, expected to get three active processes")
				}
			})

			t.Run("this must be an unsuccessful creation of a duplicate (duplicate) process", func(t *testing.T) {
				if err := workstation.PerformAsync("process_two", Payload{"id": 20, "name": "Duplicate process two"}); err == nil {
					t.Fatal("Failed, expected to get a creation error")
				} else {
					if errors.As(err, &ErrorAsyncProcessAlreadyExists) == false {
						t.Fatal("Failed, expect typed error")
					}
				}
			})

			t.Run("it shoul be successful remove process", func(t *testing.T) {
				if err := workstation.RevokeAsync("process_one"); err != nil {
					t.Fatal(err)
				}

				if err := workstation.RevokeAsync("process_one"); err == nil {
					t.Fatal("Failed, expected to get a removed error")
				} else {
					if errors.As(err, &ErrorAsyncProcessNotFoundOrAlreadyCompleted) == false {
						t.Fatal("Failed, expect typed error")
					}
				}
			})
		})

		cancel()

		// for wait all closed
		<-time.After(time.Second * 2)

		t.Run("it should be successfully give empty worstation pool", func(t *testing.T) {
			if workstation.CountAsync() != 0 {
				t.Fatal("Failed, expected to get empty pool")
			}
		})

		t.Run("it should be successfully accumulated data from processes", func(t *testing.T) {
			available := []int{10, 20, 30}

			lookupIsAvailableValue := func(id int) bool {
				for _, v := range available {
					if v == id {
						return true
					}
				}

				return false
			}

			every := map[int]struct{}{}

			for _, id := range globalCollector.All() {
				if !lookupIsAvailableValue(id.(int)) {
					t.Fatal("Failed, give no valid item")
				} else {
					every[id.(int)] = struct{}{}
				}
			}

			if len(every) != 3 {
				t.Fatal("Failed, expected to get three items")
			}

			// available is required
			for _, v := range available {
				if _, ok := every[v]; !ok {
					t.Fatalf("Failed, expected to get item %d", v)
				}
			}
		})
	})
}
