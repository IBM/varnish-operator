package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbac "k8s.io/api/rbac/v1"
)

var (
	clusterRolebindingOpts = []cmp.Option{cmpopts.IgnoreFields(rbac.ClusterRoleBinding{}, sharedIgnoreMetadata...)}
)

func EqualClusterRoleBinding(found, desired *rbac.ClusterRoleBinding) bool {
	return cmp.Equal(found, desired, clusterRolebindingOpts...)
}

func DiffClusterRoleBinding(found, desired *rbac.ClusterRoleBinding) string {
	return cmp.Diff(found, desired, clusterRolebindingOpts...)
}
