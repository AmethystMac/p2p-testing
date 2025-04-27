package sys

import "sync"

var (
	processes sync.Map
)

// Wait for the Process with name as [key] 
func WaitForProcess(key string) {

	for {
		value, ok := processes.Load(key)
		if ok && value == "done" {
			break;
		}
	}
}

// Sets the status of the Process with name as [key] and status as [value] 
func SetProcessStatus(key string, value string) {

	processes.Store(key, value)
}