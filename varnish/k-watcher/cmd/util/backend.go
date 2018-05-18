package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

// Backend is the host and port of a backend
type Backend struct {
	Host string
	Port int32
}

// String represents the "{{Host}}:{{Port}}" representation of a Backend
func (b *Backend) String() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}

// NewBackend takes a "{{Host}}:{{Port}}" representation of a Backend as initialization
func NewBackend(s string) (*Backend, error) {
	arr := strings.SplitN(s, ":", 2)
	if len(arr) != 2 {
		return nil, errors.Errorf("string not correctly formatted: %s", s)
	}
	ip := arr[0]
	portStr := arr[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.Annotate(err, "port of wrong format")
	}
	b := Backend{
		Host: ip,
		Port: int32(port),
	}
	return &b, nil
}
