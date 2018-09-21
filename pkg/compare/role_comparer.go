package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
)

var (
	roleOpts = []cmp.Option{cmpopts.IgnoreFields(rbacv1beta1.Role{}, sharedIgnoreMetadata...)}
)

func EqualRole(found, desired *rbacv1beta1.Role) bool {
	return cmp.Equal(found, desired, roleOpts...)
}

func DiffRole(found, desired *rbacv1beta1.Role) string {
	return cmp.Diff(found, desired, roleOpts...)
}
