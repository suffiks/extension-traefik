package controllers

import netv1 "k8s.io/api/networking/v1"

// +kubebuilder:validation:Pattern=`^[\w\.\-]+$`
type Host string

// +kubebuilder:validation:Pattern=`^\/.*$`
type Path string

func (p Path) IngressPath(svc, portName string) netv1.HTTPIngressPath {
	return netv1.HTTPIngressPath{
		Path: string(p),
		Backend: netv1.IngressBackend{
			Service: &netv1.IngressServiceBackend{
				Name: svc,
				Port: netv1.ServiceBackendPort{
					Name: portName,
				},
			},
		},
	}
}

type Ingress struct {
	Host  Host   `json:"host"`
	Paths []Path `json:"paths,omitempty"`
}

func (i *Ingress) Rule(svc string) netv1.IngressRule {
	ir := netv1.IngressRule{
		Host: string(i.Host),
		IngressRuleValue: netv1.IngressRuleValue{
			HTTP: &netv1.HTTPIngressRuleValue{
				Paths: []netv1.HTTPIngressPath{},
			},
		},
	}

	for _, path := range i.Paths {
		ir.HTTP.Paths = append(ir.HTTP.Paths, path.IngressPath(svc, "http"))
	}

	if len(i.Paths) == 0 {
		ir.HTTP.Paths = []netv1.HTTPIngressPath{
			Path("/").IngressPath(svc, "http"),
		}
	}

	return ir
}

// +suffiks:extension:Targets=Application,Validation=true
type Traefik struct {
	Ingresses []Ingress `json:"ingresses,omitempty"`
}
