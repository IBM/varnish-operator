package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
)

var (
	rolebindingOpts = []cmp.Option{cmpopts.IgnoreFields(rbacv1beta1.RoleBinding{}, sharedIgnoreMetadata...)}
)

func EqualRoleBinding(found, desired *rbacv1beta1.RoleBinding) bool {
	return cmp.Equal(found, desired, rolebindingOpts...)
}

func DiffRoleBinding(found, desired *rbacv1beta1.RoleBinding) string {
	return cmp.Diff(found, desired, rolebindingOpts...)
}
