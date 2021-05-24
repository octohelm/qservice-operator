package controllers

import (
	"context"
	"encoding/hex"
	"net"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/octohelm/qservice-operator/pkg/apis/serving/v1alpha1"
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	"istio.io/api/networking/v1alpha3"
	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// QEgressReconciler reconciles a QEgress object
type QEgressReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *QEgressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if !controllerutil.IsResourceRegistered(r.Client, istiov1alpha3.SchemeGroupVersion.WithKind("ServiceEntry")) {
		return nil
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.QEgress{}).
		Owns(&istiov1alpha3.ServiceEntry{}).
		Complete(r)
}

func (r *QEgressReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("namespace", request.Namespace, "name", request.Name)

	qegress := &v1alpha1.QEgress{}
	if err := r.Client.Get(ctx, request.NamespacedName, qegress); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ctx = controllerutil.ContextWithControllerClient(ctx, r.Client)

	ger := toExternalServiceEntity(qegress)
	if ger != nil {
		if err := applyServiceEntry(ctx, ger); err != nil {
			log.Error(err, "apply service entry failed")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func toExternalServiceEntity(qeg *v1alpha1.QEgress) *istiov1alpha3.ServiceEntry {
	se := &istiov1alpha3.ServiceEntry{}
	se.Namespace = qeg.Namespace

	if qeg.Spec.Egress.IP == nil && len(strings.Split(qeg.Spec.Egress.Hostname, ".")) <= 2 {
		// in cluster
		return nil
	}

	portNumber := qeg.Spec.Egress.Port

	protocol := ""
	prefix := "ext-"

	scheme := strings.ToLower(qeg.Spec.Egress.Scheme)

	switch scheme {
	case "http":
		protocol = "HTTP"
		if portNumber == 0 {
			portNumber = 80
		}
	case "https":
		protocol = "HTTPS"
		if portNumber == 0 {
			portNumber = 443
		}
	case "grpc":
		protocol = "GRPC"
	default:
		protocol = "TCP"
		prefix = "ext-" + scheme + "-"
	}

	ort := v1alpha3.Port{
		Name:     strings.ToLower(protocol) + "-" + strconv.FormatUint(uint64(portNumber), 10),
		Number:   uint32(portNumber),
		Protocol: protocol,
	}

	se.Spec.Location = v1alpha3.ServiceEntry_MESH_EXTERNAL
	se.Spec.ExportTo = []string{"."}

	if qeg.Spec.Egress.IP == nil {
		se.Name = prefix + ort.Name + "--" + qeg.Spec.Egress.Hostname
		se.Spec.Hosts = []string{qeg.Spec.Egress.Hostname}
		se.Spec.Resolution = v1alpha3.ServiceEntry_DNS
	} else {
		hip := hexIP(qeg.Spec.Egress.IP)

		se.Name = prefix + ort.Name + "-" + hip
		se.Spec.Hosts = []string{prefix + ort.Name + "-" + hip}
		se.Spec.Addresses = []string{qeg.Spec.Egress.IP.String() + "/32"}
	}

	se.Spec.Ports = []*v1alpha3.Port{&ort}

	return se
}

func hexIP(ip net.IP) string {
	v4 := ip.To4()
	return hex.EncodeToString([]byte{v4[0], v4[1], v4[2], v4[3]})
}
