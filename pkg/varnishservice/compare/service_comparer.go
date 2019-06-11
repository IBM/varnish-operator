package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1 "k8s.io/api/core/v1"
)

var (
	serviceOpts = []cmp.Option{cmpopts.IgnoreFields(v1.Service{}, sharedIgnoreMetadata...), cmpopts.IgnoreFields(v1.Service{}, sharedIgnoreStatus...)}
)

// EqualService compares 2 services for equality
func EqualService(found, desired *v1.Service) bool {
	return cmp.Equal(found, desired, serviceOpts...)
}

// DiffService generates a patch diff between 2 services
func DiffService(found, desired *v1.Service) string {
	return cmp.Diff(found, desired, serviceOpts...)
}
