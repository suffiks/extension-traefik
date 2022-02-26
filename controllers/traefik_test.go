package controllers

import (
	"os"
	"testing"

	"github.com/suffiks/suffiks/base"
	"github.com/suffiks/suffiks/extension"
	"github.com/suffiks/suffiks/extension/testutil"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefik/v1alpha1"
	faketraefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned/typed/traefik/v1alpha1/fake"
	treafikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/fake"
)

func TestTraefik(t *testing.T) {
	app := &base.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-app",
			Namespace: "mynamespace",
		},
		Spec: testutil.AppSpec(map[string]any{
			"ingresses": []map[string]any{
				{
					"host": "mydomain.org",
					"paths": []string{
						"/some/path",
					},
				},
			},
		}),
	}

	tests := []testutil.TestCase{
		testutil.SyncTest{
			Name:   "create application",
			Object: app,
			Expected: []runtime.Object{
				&treafikv1alpha1.IngressRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "some-app",
						Namespace: "mynamespace",
						OwnerReferences: []metav1.OwnerReference{
							testutil.OwnerReference("Application", "some-app"),
						},
					},
					Spec: treafikv1alpha1.IngressRouteSpec{
						Routes: []treafikv1alpha1.Route{
							{
								Match: "Host(`mydomain.org`) && PathPrefix(`/some/path`)",
								Kind:  "Rule",
								Services: []treafikv1alpha1.Service{
									{
										LoadBalancerSpec: treafikv1alpha1.LoadBalancerSpec{
											Name:      "some-app",
											Namespace: "mynamespace",
										},
									},
								},
							},
						},
						EntryPoints: []string{"web"},
					},
				},
			},
		},

		testutil.DeleteTest{
			Name:   "delete application",
			Object: app,
			ExpectedDeleted: []testutil.Deleted{
				{
					Namespace: app.Namespace,
					Name:      app.Name,
					Resource:  "ingressroutes",
				},
			},
		},
	}

	f, err := os.OpenFile("../config/crd/traefik.yaml", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	it := testutil.NewIntegrationTester(f, newExtension)
	it.Run(t, tests...)
}

func newExtension(client *fake.Clientset) extension.Extension[*Traefik] {
	return &TraefikExtension{
		Traefik: &fakeWrap{client: client},
	}
}

type fakeWrap struct {
	client *fake.Clientset
}

func (f *fakeWrap) Discovery() discovery.DiscoveryInterface {
	return f.Discovery()
}

func (f *fakeWrap) TraefikV1alpha1() traefikv1alpha1.TraefikV1alpha1Interface {
	return &faketraefikv1alpha1.FakeTraefikV1alpha1{
		Fake: &f.client.Fake,
	}
}
