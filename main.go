package main

import (
	"flag"
	"os"

	servingapis "github.com/octohelm/qservice-operator/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/controllers"
	istioapis "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(istioapis.AddToScheme(scheme))

	utilruntime.Must(servingapis.AddToScheme(scheme))
}

func main() {
	ctrlOpt := ctrl.Options{
		Scheme:           scheme,
		Port:             9443,
		LeaderElectionID: "74b83f88.octohelm.tech",
	}

	flag.StringVar(&ctrlOpt.Namespace, "watch-namespace", os.Getenv("WATCH_NAMESPACE"), "watch namespace")
	flag.StringVar(&ctrlOpt.MetricsBindAddress, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&ctrlOpt.LeaderElection, "enable-leader-election", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrlOpt)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = controllers.AddToManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	setupLog.Info("starting manager")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
