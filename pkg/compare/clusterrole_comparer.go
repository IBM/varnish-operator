package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
)

var (
	clusterRoleOpts = []cmp.Option{cmpopts.IgnoreFields(rbacv1beta1.ClusterRole{}, sharedIgnoreMetadata...)}
)

func EqualClusterRole(found, desired *rbacv1beta1.ClusterRole) bool {
	return cmp.Equal(found, desired, clusterRoleOpts...)
}

func DiffClusterRole(found, desired *rbacv1beta1.ClusterRole) string {
	return cmp.Diff(found, desired, clusterRoleOpts...)
}
