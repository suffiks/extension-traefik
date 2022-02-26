package controllers

import (
	"context"

	"github.com/suffiks/suffiks/extension"
	traefikcrd "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned"
	treafikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	_ extension.Extension[*Traefik]            = &TraefikExtension{}
	_ extension.ValidatableExtension[*Traefik] = &TraefikExtension{}
)

//+kubebuilder:rbac:groups=traefik.containo.us,resources=ingressroutes,verbs=get;list;watch;create;update;patch;delete

type TraefikExtension struct {
	Traefik traefikcrd.Interface

	AllowedDomains []string
}

func (t *TraefikExtension) Sync(ctx context.Context, owner extension.Owner, obj *Traefik, rw *extension.ResponseWriter) error {
	irclient := t.Traefik.TraefikV1alpha1().IngressRoutes(owner.Namespace())
	exists := true

	ir, err := irclient.Get(ctx, owner.Name(), metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

		ir = &treafikv1alpha1.IngressRoute{
			ObjectMeta: owner.ObjectMeta(),
			Spec: treafikv1alpha1.IngressRouteSpec{
				EntryPoints: []string{"web"},
			},
		}
		exists = false
	} else {
		ir.Spec.Routes = nil
	}

	for _, rule := range obj.Ingresses {
		for _, path := range rule.Paths {
			ir.Spec.Routes = append(ir.Spec.Routes, treafikv1alpha1.Route{
				Match: "Host(`" + string(rule.Host) + "`) && PathPrefix(`" + string(path) + "`)",
				Kind:  "Rule",
				Services: []treafikv1alpha1.Service{
					{
						LoadBalancerSpec: treafikv1alpha1.LoadBalancerSpec{
							Name:      owner.Name(),
							Namespace: owner.Namespace(),
						},
					},
				},
			})
		}
	}

	if exists {
		_, err = irclient.Update(ctx, ir, metav1.UpdateOptions{})
	} else {
		_, err = irclient.Create(ctx, ir, metav1.CreateOptions{})
	}
	return err
}

func (t *TraefikExtension) Delete(ctx context.Context, owner extension.Owner, obj *Traefik) error {
	irclient := t.Traefik.TraefikV1alpha1().IngressRoutes(owner.Namespace())
	err := irclient.Delete(ctx, owner.Name(), metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (t *TraefikExtension) Validate(ctx context.Context, typ extension.ValidationType, owner extension.Owner, newObj, oldObj *Traefik) ([]extension.ValidationErrors, error) {
	if typ == extension.ValidationDelete {
		return nil, nil
	}

	if len(t.AllowedDomains) == 0 {
		return nil, nil
	}

	var errs []extension.ValidationErrors

	return errs, nil
}
