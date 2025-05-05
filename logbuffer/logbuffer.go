package logbuffer

import "sync"

var logs []string
var mu sync.Mutex

func AddLog(entry string) {
	mu.Lock()
	logs = append(logs, entry)
	mu.Unlock()
}

func GetLogs() []string {
	mu.Lock()
	defer mu.Unlock()
	return append([]string(nil), logs...) // return copy
}
