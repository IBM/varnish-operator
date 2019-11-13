package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbac "k8s.io/api/rbac/v1"
)

var (
	roleOpts = []cmp.Option{cmpopts.IgnoreFields(rbac.Role{}, sharedIgnoreMetadata...)}
)

func EqualRole(found, desired *rbac.Role) bool {
	return cmp.Equal(found, desired, roleOpts...)
}

func DiffRole(found, desired *rbac.Role) string {
	return cmp.Diff(found, desired, roleOpts...)
}
