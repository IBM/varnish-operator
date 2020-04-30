package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1 "k8s.io/api/core/v1"
)

var (
	configMapOpts = []cmp.Option{cmpopts.IgnoreFields(v1.ConfigMap{}, sharedIgnoreMetadata...)}
)

func EqualConfigMap(found, desired *v1.ConfigMap) bool {
	return cmp.Equal(found, desired, configMapOpts...)
}

func DiffConfigMap(found, desired *v1.ConfigMap) string {
	return cmp.Diff(found, desired, append(configMapOpts, cmpopts.IgnoreFields(v1.ConfigMap{}, "Data"))...)
}
