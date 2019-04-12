package controller

import (
	"fmt"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"regexp"
	"sort"
	"strings"
)

var (
	varnishArgsKeyRegexp = regexp.MustCompile("^-\\w")
	disallowedArgs       = [][]string{
		{"-n"},
	}
)

func getSanitizedVarnishArgs(spec *icmapiv1alpha1.VarnishServiceSpec) []string {
	varnishArgsOverrides := [][]string{
		{"-F"},
		{"-a", fmt.Sprintf("0.0.0.0:%d", icmapiv1alpha1.VarnishPort)},
		{"-S", "/etc/varnish/secret"},
		{"-f", "/etc/varnish/" + spec.VCLConfigMap.EntrypointFile},
		{"-T", fmt.Sprintf("127.0.0.1:%d", icmapiv1alpha1.VarnishAdminPort)},
	}

	rawArgs := spec.Deployment.Container.VarnishArgs
	var parsedArgs [][]string

	// parse arguments and remove ones that should be overridden
	argsCount := len(rawArgs)
	for i := 0; i < argsCount; {
		var nextArg []string

		// add arg key to output
		nextArg = append(nextArg, rawArgs[i])
		i++
		// if there is an arg value (as defined by NOT being a key), add it to output
		if i < argsCount && !varnishArgsKeyRegexp.MatchString(rawArgs[i]) {
			nextArg = append(nextArg, rawArgs[i])
			i++
		}

		// skip the argument if it's present among overrides or not allowed to specify at all
		if argSpecified(varnishArgsOverrides, nextArg[0]) || argSpecified(disallowedArgs, nextArg[0]) {
			continue
		}

		parsedArgs = append(parsedArgs, nextArg)
	}

	varnishArgs := append(parsedArgs, varnishArgsOverrides...)

	// sort the arguments so they won't appear in different order in different reconcile loops and trigger redeployment
	sort.SliceStable(varnishArgs, func(i, j int) bool {
		// we can't just compare by keys as there can be multiple equal keys that set different parameters
		// e.g. "-p default_ttl=200 -p default_grace=1000". So compare by combiantion of key and value
		return strings.Join(varnishArgs[i], " ") < strings.Join(varnishArgs[j], " ")
	})

	var sanitizedArgs []string
	for _, value := range varnishArgs {
		sanitizedArgs = append(sanitizedArgs, value...)
	}

	return sanitizedArgs
}

// argSpecified checks if the user specified the argument
func argSpecified(args [][]string, arg string) bool {
	for _, value := range args {
		if value[0] == arg {
			return true
		}
	}
	return false

}
