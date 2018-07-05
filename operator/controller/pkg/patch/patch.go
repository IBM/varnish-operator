package patch

import (
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/mergepatch"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

var json = jsoniter.ConfigFastest

type noChange struct {
	errors.Err
}

// NoChangef creates a new error of type NoChange
func NoChangef(format string, args ...interface{}) error {
	err := errors.NewErr(format+" no change", args)
	err.SetLocation(1)
	return &noChange{Err: err}
}

// IsNoChange checks if an error's root cause was of type NoChange
func IsNoChange(err error) bool {
	_, ok := errors.Cause(err).(*noChange)
	return ok
}

// NewStrategicMergePatch generates the required patch bytes for the Kubernetes API
func NewStrategicMergePatch(lastModified, desired, current runtime.Object, fns ...mergepatch.PreconditionFunc) ([]byte, error) {
	lastModifiedJSON, err := json.Marshal(lastModified)
	if err != nil {
		return nil, errors.Annotate(err, "could not marshall last modified state")
	}
	desiredJSON, err := json.Marshal(desired)
	if err != nil {
		return nil, errors.Annotate(err, "could not marshall desired state")
	}
	currentJSON, err := json.Marshal(current)
	if err != nil {
		return nil, errors.Annotate(err, "could not marshal current state")
	}

	lookupPatchMeta, err := strategicpatch.NewPatchMetaFromStruct(current)
	if err != nil {
		return nil, errors.Annotate(err, "problem generating three-way patch")
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch(lastModifiedJSON, desiredJSON, currentJSON, lookupPatchMeta, true, fns...)
	if len(patch) == 0 || string(patch) == "{}" {
		return nil, NoChangef("object patch is empty")
	}
	return patch, nil
}