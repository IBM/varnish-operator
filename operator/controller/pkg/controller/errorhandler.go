package controller

import (
	"github.com/juju/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func handleError(err error) {
	stackTrace := errors.Details(err)
	errWithStack := errors.NewErr(stackTrace)
	errWithStack.SetLocation(1)
	utilruntime.HandleError(&errWithStack)
}
