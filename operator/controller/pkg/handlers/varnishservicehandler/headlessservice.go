package varnishservicehandler

import (
	icmapiv1alpha1 "icm-varnish-k8s-operator/operator/controller/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/operator/controller/pkg/config"

	"fmt"
	"icm-varnish-k8s-operator/operator/controller/pkg/patch"

	"github.com/json-iterator/go"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

var json = jsoniter.ConfigFastest

var LastAppliedConfigurationKey = "icm.ibm.com/last-applied-configuration"

func applyHeadlessService(client *kubernetes.Clientset, globalConf *config.Config, vs *icmapiv1alpha1.VarnishService) error {
	serviceClient := client.CoreV1().Services(vs.Namespace)

	var lastAppliedState, desiredState, currState *v1.Service
	var err error

	currState, err = serviceClient.Get(vs.Name, metav1.GetOptions{})
	if kerrors.IsNotFound(err) {
		currState = nil
	} else if err != nil {
		return errors.Annotate(err, "could not apply headless service")
	}

	varnishBackedPort, otherPorts, err := extractVarnishPort(vs.Spec.Service.Spec.Ports, vs.Spec.Service.Spec.VarnishPortName)
	if err != nil {
		return errors.Trace(err)
	}

	conf := headlessConfig{
		ServiceName:       fmt.Sprintf("%s-headless-service", vs.Name),
		AppLabels:         globalConf.VarnishCommonLabels,
		AppSelectors:      vs.Spec.Service.Spec.Selector,
		VarnishBackedPort: *varnishBackedPort,
		OtherPorts:        otherPorts,
	}
	desiredState, err = newHeadlessService(globalConf, &conf)
	if err != nil {
		return errors.Annotate(err, "could not create new headless service")
	}

	lastAppliedState, err = lastAppliedHeadlessService(currState)
	if err != nil {
		return errors.Annotate(err, "could not retrieve last applied state")
	}

	patchData, err := patch.NewStrategicMergePatch(lastAppliedState, desiredState, currState)
	if err != nil {
		return errors.Trace(err)
	}
	res, err := serviceClient.Patch(vs.Name, types.StrategicMergePatchType, patchData)
	if err != nil {
		return errors.Annotate(err, "problem executing patch")
	}

	log.WithField("headless-service", res).Info("applied headless service")
	return nil
}

func lastAppliedHeadlessService(currState *v1.Service) (*v1.Service, error) {
	if currState == nil {
		return nil, nil
	}

	lastAppliedJson, found := currState.Annotations[LastAppliedConfigurationKey]
	if !found {
		return nil, kerrors.NewNotFound(v1.Resource("service"), currState.Name)
	}

	s := v1.Service{}
	if err := json.UnmarshalFromString(lastAppliedJson, &s); err != nil {
		return nil, errors.Annotate(err, "could not unmarshal last-applied-configuration json")
	}
	return &s, nil
}

func extractVarnishPort(ports []v1.ServicePort, varnishPortName string) (*v1.ServicePort, []v1.ServicePort, error) {
	var varnishBackedPort *v1.ServicePort
	otherPorts := make([]v1.ServicePort, len(ports)-1)

	for _, port := range ports {
		if port.Name == varnishPortName {
			if varnishBackedPort != nil {
				return nil, nil, errors.New("more than one port had name of VarnishBackedPort")
			}
			varnishBackedPort = &port
		} else {
			otherPorts = append(otherPorts, port)
		}
	}
	return varnishBackedPort, otherPorts, nil
}

func deleteHeadlessService(client *kubernetes.Clientset, vs *icmapiv1alpha1.VarnishService) error {
	serviceClient := client.CoreV1().Services(vs.Namespace)

	return serviceClient.Delete(vs.Name, &metav1.DeleteOptions{})
}

type headlessConfig struct {
	ServiceName       string
	AppLabels         map[string]string
	AppSelectors      map[string]string
	VarnishBackedPort v1.ServicePort
	OtherPorts        []v1.ServicePort
}

func newHeadlessService(globalConf *config.Config, headlessConf *headlessConfig) (*v1.Service, error) {
	s := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   headlessConf.ServiceName,
			Labels: headlessConf.AppLabels,
			Annotations: map[string]string{
				"icm.ibm.com/owner":               globalConf.VarnishName,
				"icm.ibm.com/varnish-backed-port": string(headlessConf.VarnishBackedPort.Port),
			},
		},
		Spec: v1.ServiceSpec{
			Ports:     append(headlessConf.OtherPorts, headlessConf.VarnishBackedPort),
			Selector:  headlessConf.AppSelectors,
			ClusterIP: "None",
			Type:      v1.ServiceTypeClusterIP,
		},
	}
	var err error
	s.Annotations[LastAppliedConfigurationKey], err = json.MarshalToString(s)

	if err != nil {
		return nil, errors.Annotate(err, "could not marshal Service to JSON")
	}
	return &s, nil
}
