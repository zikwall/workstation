package main

import (
	"fmt"
	"log"
)

func infoSuccessKillprocess(id string, pid int) {
	log.Println(fmt.Sprintf("[INFO] Async task #%s sub process PID %d successfully killed", id, pid))
}

func errorFailedkillprocess(id string, pid int, err error) {
	log.Println(
		fmt.Sprintf("[ERROR] Failed to kill async task #%s sub process PID %d %s", id, pid, err),
	)
}

func errorProcessKilled(id string, pid int, err error) {
	log.Println(
		fmt.Sprintf(`[ERROR] The async task #%s process PID %d was terminated with an error, the task is removed from the pool and will be restarted in the future. Previous error '%s'`, id, pid, err),
	)
}

func errorCloseProcessStdoutReader(id string, err error) {
	log.Println(fmt.Sprintf("[ERROR] Reading from stdout for asynctask #%s completed (with error), cause %s", id, err))
}

func infoCloseProcessStdoutReader(id string) {
	log.Println(fmt.Sprintf("[INFO] Read from stdout for async task #%s completed, successfully!", id))
}
