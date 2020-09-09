package controllers

import (
	"github.com/octohelm/qservice-operator/controllers/deployment"
	"github.com/octohelm/qservice-operator/controllers/qservice"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func AddToManager(mgr ctrl.Manager) error {
	return SetupReconcilerWithManager(
		mgr,
		&qservice.QServiceReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("QService"),
			Scheme: mgr.GetScheme(),
		},
		&deployment.DeploymentReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Deployment"),
			Scheme: mgr.GetScheme(),
		},
	)
}

type Reconciler interface {
	SetupWithManager(mgr ctrl.Manager) error
}

func SetupReconcilerWithManager(mgr manager.Manager, reconcilers ...Reconciler) error {
	for i := range reconcilers {
		if err := reconcilers[i].SetupWithManager(mgr); err != nil {
			return err
		}
	}
	return nil
}
