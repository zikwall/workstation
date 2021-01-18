[![build](https://github.com/zikwall/workstation/workflows/tests/badge.svg)](https://github.com/zikwall/workstation/actions)

<div align="center">
  <h1>Workstation</h1>
  <h5>Simple, Powerful and productive asynchronous task manager written in Go</h5>
</div>

#### Types

- [x] Workstation - suitable for long-running processes
- [ ] WorkstationWorkerPool - pool of tasks that are processed in a limited number of threads

#### How to use

- `go get -u github.com/zikwall/workstation`

##### Workstation. Define custom worker

```go
package main

import (
	"context"
	"github.com/zikwall/workstation"
)

type MyWorker struct{}

func (w *MyWorker) Perform(ctx context.Context, key string, payload workstation.Payload) {
	//.. <your custom code here>
	//.. anything longer
	//.. or unit process
}

// and create workstation ...
```

##### Workstation. Create workstation

```go
ctx, cancel := context.WithCancel(context.Background())

station := workstation.BuildWorkstation(ctx, &MyWorker{})

// run sub-processes

err := station.PerformAsync(
    "first_process", workstation.Payload{"a": 1, "b": "6", "c": sampleFunctionC},
)

// .. <another sub-processes>
// .. <your custom code>

cancel()

if err := station.Shutdown(); err != nil {
	log.Fatal(err)
}
```

#### Todo

- [ ] **Workstation**: configuration
- [ ] **Workstation**: logs
- [ ] **Workstation**: limitation of parallel processes
- [ ] **Workstation**: max retry the process if the process failed with an error
- [ ] **Workstation**: server for pull/push jobs
- [ ] **Workstation**: web interface with auth, dashboards, control panel
- [ ] **Workstation**: process states: processed, busy, failed, dead (if retry > N), cancelled
- [ ] **WorkstationWorkerPool**: need to implement