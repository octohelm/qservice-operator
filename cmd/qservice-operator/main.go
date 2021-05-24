package main

import (
	"flag"
	"os"

	"github.com/octohelm/qservice-operator/pkg/controllerutil"

	"github.com/octohelm/qservice-operator/internal/version"
	"github.com/octohelm/qservice-operator/pkg/apis/serving"
	servingapis "github.com/octohelm/qservice-operator/pkg/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/controllers"
	"github.com/pkg/errors"
	istioapis "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(istioapis.AddToScheme(scheme))
	utilruntime.Must(servingapis.AddToScheme(scheme))
}

func start(ctrlOpt ctrl.Options) error {
	restConfig := ctrl.GetConfigOrDie()

	if err := controllerutil.ApplyCRDs(restConfig, serving.CRDs...); err != nil {
		return errors.Wrap(err, "unable to create crds")
	} else {
		ctrl.Log.WithName("crd").Info("crds created")
	}

	mgr, err := ctrl.NewManager(restConfig, ctrlOpt)
	if err != nil {
		return errors.Wrap(err, "unable to start manager")
	}

	if err := controllers.SetupWithManager(mgr); err != nil {
		return errors.Wrap(err, "unable to create controller")
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return errors.Wrap(err, "problem running manager")
	}
	return nil
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(false)))

	ctrlOpt := ctrl.Options{
		Scheme:           scheme,
		Port:             9443,
		LeaderElectionID: "74b83f88.octohelm.tech",
		Logger:           ctrl.Log.WithValues("qservice-operator", version.Version),
	}

	flag.StringVar(&ctrlOpt.Namespace, "watch-namespace", os.Getenv("WATCH_NAMESPACE"), "watch namespace")
	flag.StringVar(&ctrlOpt.MetricsBindAddress, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&ctrlOpt.LeaderElection, "enable-leader-election", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	flag.Parse()

	if err := start(ctrlOpt); err != nil {
		ctrl.Log.WithName("setup").Error(err, "")
	}
}
