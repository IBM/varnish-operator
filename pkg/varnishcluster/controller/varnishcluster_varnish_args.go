package controller

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	vcapi "github.com/cin/varnish-operator/api/v1alpha1"
)

var (
	varnishArgsKeyRegexp = regexp.MustCompile(`^-\\w`)
	disallowedArgs       = [][]string{
		{"-n"},
		{"-f"},
	}
)

func getSanitizedVarnishArgs(spec *vcapi.VarnishClusterSpec) []string {
	varnishArgsOverrides := [][]string{
		{"-F"},
		{"-S", "/etc/varnish-secret/secret"},
		{"-b", "127.0.0.1:0"}, //start a varnishd without predefined backend. It has to be overridden by settings from ConfigMap
		{"-T", fmt.Sprintf("0.0.0.0:%d", vcapi.VarnishAdminPort)},
	}

	rawArgs := spec.Varnish.Args
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
	varnishArgs = append(varnishArgs, []string{"-a", fmt.Sprintf("0.0.0.0:%d", vcapi.VarnishPort)})

	// sort the arguments so they won't appear in different order in different reconcile loops and trigger redeployment
	sort.SliceStable(varnishArgs, func(i, j int) bool {
		// we can't just compare by keys as there can be multiple equal keys that set different parameters
		// e.g. "-p default_ttl=200 -p default_grace=1000". So compare by combination of key and value
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
