package controller

import (
	"icm-varnish-k8s-operator/pkg/apis"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishservice/config"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"

	"go.uber.org/zap"

	icmv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"

	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}
var stsKey = types.NamespacedName{Name: "foo-varnish-statefulset", Namespace: "default"}

const timeout = time.Second * 5

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &icmv1alpha1.VarnishService{
		ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"},
		Spec: icmv1alpha1.VarnishServiceSpec{
			Service: icmv1alpha1.VarnishServiceService{
				VarnishPort: v1.ServicePort{
					Port: 80,
				},
			},
			VCLConfigMap: icmv1alpha1.VarnishVCLConfigMap{
				Name:           "vcl-files",
				EntrypointFile: "default.vcl",
			},
		},
	}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	//zapLogr, err := zap.NewDevelopment() //uncomment if you want to see operator logs for debugging
	zapLogr := zap.NewNop()
	g.Expect(err).NotTo(gomega.HaveOccurred())
	testCfg := &config.Config{
		CoupledVarnishImage: "us.icr.io/icm-varnish/varnish:0.18.0",
	}

	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		g.Expect(err).NotTo(gomega.HaveOccurred())
	}

	testLogger := &logger.Logger{
		SugaredLogger: zapLogr.Sugar(),
	}

	reconciler, requests := SetupTestReconcile(NewVarnishReconciler(mgr, testCfg, testLogger))
	g.Expect(Add(reconciler, mgr, testLogger)).NotTo(gomega.HaveOccurred())
	defer close(StartTestManager(mgr, g))

	// Create the VarnishService object and expect the Reconcile and StatefulSet to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer func() {
		err := c.Delete(context.TODO(), instance)
		if err != nil {
			t.Logf("failed to delete the instance. Error: %v", err)
		}
	}()

	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	statefulSet := &appsv1.StatefulSet{}
	g.Eventually(func() error { return c.Get(context.TODO(), stsKey, statefulSet) }, timeout).
		Should(gomega.Succeed())

	// Delete the StatefulSet and expect Reconcile to be called for StatefulSet deletion
	g.Expect(c.Delete(context.TODO(), statefulSet)).NotTo(gomega.HaveOccurred())
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	g.Eventually(func() error { return c.Get(context.TODO(), stsKey, statefulSet) }, timeout).
		Should(gomega.Succeed())

	// Manually delete StatefulSet since GC isn't enabled in the test control plane
	g.Expect(c.Delete(context.TODO(), statefulSet)).To(gomega.Succeed())
}
