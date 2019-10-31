package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbac "k8s.io/api/rbac/v1"
)

var (
	clusterRoleOpts = []cmp.Option{cmpopts.IgnoreFields(rbac.ClusterRole{}, sharedIgnoreMetadata...)}
)

func EqualClusterRole(found, desired *rbac.ClusterRole) bool {
	return cmp.Equal(found, desired, clusterRoleOpts...)
}

func DiffClusterRole(found, desired *rbac.ClusterRole) string {
	return cmp.Diff(found, desired, clusterRoleOpts...)
}
