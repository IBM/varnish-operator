/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishcluster/config"
	vcreconcile "icm-varnish-k8s-operator/pkg/varnishcluster/reconcile"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sync"
	"testing"

	"github.com/go-logr/zapr"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config                            //config for the rest client
var k8sClient client.Client                     //k8s client that will use the config above to point to the test environment
var testEnv *envtest.Environment                //brings up the control plane that you can connect to using the client above
var requestsChan = make(chan reconcile.Request) //receives a value every time a reconcile loop finishes
var mgrStopCh = make(chan struct{})             //stops the manager by sending a value to the channel
var waitGroup = &sync.WaitGroup{}               //waits until the reconcile loops finish. Used to gracefully shutdown the environment
var reconcileChan = make(chan event.GenericEvent) //can be used to send manually reconcile events. Useful for testing.

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{envtest.NewlineReporter{}},
	)
}

var _ = BeforeSuite(func(done Done) {
	//logr, destWriter := logger.NewLogger("console", zapcore.DebugLevel), GinkgoWriter //uncomment and replace with the line below to have logging
	logr, destWriter := logger.NewNopLogger(), GinkgoWriter
	ctrl.SetLogger(zapr.NewLogger(logr.Desugar()))
	logf.SetLogger(zapr.NewLogger(logr.Desugar()))

	logf.SetLogger(zap.New(func(o *zap.Options) {
		o.DestWritter = destWriter
	}))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = icmv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(mgr).ToNot(BeNil())

	vcCtrl := &ReconcileVarnishCluster{
		logger:             logr,
		scheme:             scheme.Scheme,
		Client:             k8sClient,
		events:             NewEventHandler(&record.FakeRecorder{Events: make(chan string)}),
		config:             &config.Config{CoupledVarnishImage: "varnish-operator:0.20.1"},
		reconcileTriggerer: vcreconcile.NewReconcileTriggerer(logr, reconcileChan),
	}

	var testReconciler reconcile.Reconciler
	testReconciler, requestsChan = SetupTestReconcile(vcCtrl)

	err = SetupVarnishReconciler(testReconciler, mgr, reconcileChan)
	Expect(err).ToNot(HaveOccurred())

	mgrStopCh = StartTestManager(mgr)
	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	close(mgrStopCh)      //tell the manager to shutdown
	waitGroup.Wait()      //wait for all reconcile loops to be finished
	err := testEnv.Stop() //stop the test control plane (etcd, kube-apiserver)
	Expect(err).ToNot(HaveOccurred())
})

func SetupTestReconcile(inner reconcile.Reconciler) (reconcile.Reconciler, chan reconcile.Request) {
	requests := make(chan reconcile.Request)
	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		waitGroup.Add(1)
		result, err := inner.Reconcile(req)
		waitGroup.Done()
		requests <- req
		return result, err
	})
	return fn, requests
}

func StartTestManager(mgr manager.Manager) chan struct{} {
	stop := make(chan struct{})
	go func() {
		Expect(mgr.Start(stop)).NotTo(HaveOccurred())
	}()
	return stop
}

// As the test control plane doesn't support garbage collection, this function is used to clean up resources
// Designed to not fail if the resource is not found
func CleanUpCreatedResources(vcName, vcNamespace string) {
	err := k8sClient.DeleteAllOf(context.Background(), &icmv1alpha1.VarnishCluster{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
	err = k8sClient.DeleteAllOf(context.Background(), &rbac.ClusterRole{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
	err = k8sClient.DeleteAllOf(context.Background(), &rbac.ClusterRoleBinding{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
	err = k8sClient.DeleteAllOf(context.Background(), &v1.ConfigMap{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
	haveNoErrorOrNotFoundError := Or(BeNil(), BeAssignableToTypeOf(&errors.StatusError{}))
	err = k8sClient.Delete(context.Background(), &v1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: vcNamespace, Name: vcName + "-headless-service"}})
	Expect(err).To(haveNoErrorOrNotFoundError)
	err = k8sClient.Delete(context.Background(), &v1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: vcNamespace, Name: vcName + "-no-cache"}})
	Expect(err).To(haveNoErrorOrNotFoundError)
	err = k8sClient.Delete(context.Background(), &v1.Service{ObjectMeta: metav1.ObjectMeta{Namespace: vcNamespace, Name: vcName}})
	Expect(err).To(haveNoErrorOrNotFoundError)
	err = k8sClient.DeleteAllOf(context.Background(), &policyv1beta1.PodDisruptionBudget{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
	err = k8sClient.DeleteAllOf(context.Background(), &rbac.Role{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
	err = k8sClient.DeleteAllOf(context.Background(), &rbac.RoleBinding{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
	err = k8sClient.DeleteAllOf(context.Background(), &v1.ServiceAccount{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
	err = k8sClient.DeleteAllOf(context.Background(), &apps.StatefulSet{}, client.InNamespace(vcNamespace))
	Expect(err).To(BeNil())
}
