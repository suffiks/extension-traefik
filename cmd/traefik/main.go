package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/suffiks/extension-traefik/controllers"
	"github.com/suffiks/suffiks/extension"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned"
	"k8s.io/client-go/rest"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	client, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}

	var domains []string
	if d := os.Getenv("ALLOWED_DOMAINS"); d != "" {
		domains = strings.Split(d, ",")
	}

	ext := &controllers.TraefikExtension{
		Traefik:        client,
		AllowedDomains: domains,
	}

	fmt.Println("Listening on :4269")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	return extension.Serve[*controllers.Traefik](ctx, ":4269", ext)
}
