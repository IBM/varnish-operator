package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
)

var (
	clusterRolebindingOpts = []cmp.Option{cmpopts.IgnoreFields(rbacv1beta1.ClusterRoleBinding{}, sharedIgnoreMetadata...)}
)

func EqualClusterRoleBinding(found, desired *rbacv1beta1.ClusterRoleBinding) bool {
	return cmp.Equal(found, desired, clusterRolebindingOpts...)
}

func DiffClusterRoleBinding(found, desired *rbacv1beta1.ClusterRoleBinding) string {
	return cmp.Diff(found, desired, clusterRolebindingOpts...)
}
