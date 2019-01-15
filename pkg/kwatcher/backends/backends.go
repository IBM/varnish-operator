package backends

import (
	"sync"
)

type Backend struct {
	IP         string
	NodeLabels map[string]string
	PodName    string
}

var (
	backends     []Backend
	backendsLock sync.RWMutex
)

func ReplaceBackends(newBackends []Backend) []Backend {
	backendsLock.Lock()
	defer backendsLock.Unlock()

	backends = make([]Backend, len(newBackends))
	copy(backends, newBackends)
	return backends
}

func GetCurrBackends() []Backend {
	backendsLock.RLock()
	defer backendsLock.RUnlock()

	backendsCopy := make([]Backend, len(backends))
	copy(backendsCopy, backends)
	return backendsCopy
}
