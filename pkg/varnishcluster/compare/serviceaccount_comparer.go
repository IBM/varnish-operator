package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1 "k8s.io/api/core/v1"
)

var (
	serviceAccountOpts = []cmp.Option{cmpopts.IgnoreFields(v1.ServiceAccount{}, append(sharedIgnoreMetadata, "ImagePullSecrets")...)}
)

func withServiceAccountInheritance(found, desired *v1.ServiceAccount) {
	if desired.Secrets == nil {
		desired.Secrets = found.Secrets
	}
}

func EqualServiceAccount(found, desired *v1.ServiceAccount) bool {
	desiredCopy := &v1.ServiceAccount{}
	desired.DeepCopyInto(desiredCopy)
	withServiceAccountInheritance(found, desiredCopy)
	return cmp.Equal(found, desiredCopy, serviceAccountOpts...)
}

func DiffServiceAccount(found, desired *v1.ServiceAccount) string {
	desiredCopy := &v1.ServiceAccount{}
	desired.DeepCopyInto(desiredCopy)
	withServiceAccountInheritance(found, desiredCopy)
	return cmp.Diff(found, desiredCopy, serviceAccountOpts...)
}
