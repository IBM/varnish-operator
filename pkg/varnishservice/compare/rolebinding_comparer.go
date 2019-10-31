package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbac "k8s.io/api/rbac/v1"
)

var (
	rolebindingOpts = []cmp.Option{cmpopts.IgnoreFields(rbac.RoleBinding{}, sharedIgnoreMetadata...)}
)

func EqualRoleBinding(found, desired *rbac.RoleBinding) bool {
	return cmp.Equal(found, desired, rolebindingOpts...)
}

func DiffRoleBinding(found, desired *rbac.RoleBinding) string {
	return cmp.Diff(found, desired, rolebindingOpts...)
}
