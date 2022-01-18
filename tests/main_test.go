package tests

import (
	"k8s.io/client-go/kubernetes"
	"testing"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/logger"
	"k8s.io/client-go/rest"

	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	debugLogsDir = "/tmp/debug-logs/"
)

var (
	k8sClient  client.Client
	restConfig *rest.Config
	kubeClient *kubernetes.Clientset
	tailLines  int64 = 30
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(func(o *zap.Options) { o.DestWriter = GinkgoWriter }))
	logr := logger.NewLogger("console", zapcore.DebugLevel)
	By("bootstrapping test environment")

	var err error
	err = vcapi.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	restConfig, err = ctrl.GetConfig()
	if err != nil {
		logr.Fatalf("unable to set up client config: %v", err)
	}

	k8sClient, err = client.New(restConfig, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	// Create client test. We use kubernetes package bc currently only it has GetLogs method.
	kubeClient, err = kubernetes.NewForConfig(restConfig)
	Expect(err).ToNot(HaveOccurred())
})
