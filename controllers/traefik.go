package controllers

import (
	"context"
	"fmt"

	"github.com/suffiks/suffiks/extension"
	traefikcrd "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned"
	treafikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	_ extension.Extension[*Traefik]            = &TraefikExtension{}
	_ extension.ValidatableExtension[*Traefik] = &TraefikExtension{}
)

//+kubebuilder:rbac:groups=traefik.containo.us,resources=ingressroutes,verbs=get;list;watch;create;update;patch;delete

type TraefikExtension struct {
	Traefik traefikcrd.Interface

	allowedDomains *tree
}

func (t *TraefikExtension) Sync(ctx context.Context, owner extension.Owner, obj *Traefik, rw *extension.ResponseWriter) error {
	fmt.Println("Sync")
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
							Port:      intstr.FromString("http"),
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
	fmt.Println("Delete")
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

	if t.allowedDomains == nil {
		return nil, nil
	}

	var errs []extension.ValidationErrors
	for _, rule := range newObj.Ingresses {
		if !t.allowedDomains.Contains(rule.Host) {
			errs = append(errs, extension.ValidationErrors{
				Path:   "ingresses.host",
				Value:  string(rule.Host),
				Detail: fmt.Sprintf("Host %s is not allowed", rule.Host),
			})
		}
	}

	return errs, nil
}

func (t *TraefikExtension) AddAllowedDomains(domains []string) {
	if len(domains) == 0 {
		return
	}

	t.allowedDomains = &tree{
		node: make(node),
	}
	for _, domain := range domains {
		t.allowedDomains.Add(domain)
	}
}
