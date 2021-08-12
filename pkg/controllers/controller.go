package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func SetupWithManager(mgr ctrl.Manager) error {
	return SetupReconcilerWithManager(
		mgr,
		&QServiceReconciler{
			Client: ClientWithoutCache(mgr.GetClient(), mgr.GetAPIReader()),
			Log:    mgr.GetLogger().WithName("controllers").WithName("QService"),
			Scheme: mgr.GetScheme(),
		},
		&DeploymentReconciler{
			Client: ClientWithoutCache(mgr.GetClient(), mgr.GetAPIReader()),
			Log:    mgr.GetLogger().WithName("controllers").WithName("Deployment"),
			Scheme: mgr.GetScheme(),
		},
		&IngressReconciler{
			Client: ClientWithoutCache(mgr.GetClient(), mgr.GetAPIReader()),
			Log:    mgr.GetLogger().WithName("controllers").WithName("Ingress"),
			Scheme: mgr.GetScheme(),
		},
		&ServiceReconciler{
			Client: ClientWithoutCache(mgr.GetClient(), mgr.GetAPIReader()),
			Log:    mgr.GetLogger().WithName("controllers").WithName("Service"),
			Scheme: mgr.GetScheme(),
		},
		&QIngressReconciler{
			Client: ClientWithoutCache(mgr.GetClient(), mgr.GetAPIReader()),
			Log:    mgr.GetLogger().WithName("controllers").WithName("QIngress"),
			Scheme: mgr.GetScheme(),
		},
		&QEgressReconciler{
			Client: ClientWithoutCache(mgr.GetClient(), mgr.GetAPIReader()),
			Log:    mgr.GetLogger().WithName("controllers").WithName("QEgress"),
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
