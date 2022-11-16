package sftp

import "sync"

var lock sync.RWMutex
var clientCache map[string]*Client

func init() {
	clientCache = make(map[string]*Client)
}
