package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/suffiks/extension-traefik/controllers"
	"github.com/suffiks/suffiks/extension"
	"k8s.io/client-go/rest"

	// This package requires a lot of replacements in go.mod to work.
	// Might be worth looking into another way to do this.
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config-file", "", "path to config file")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	client, err := versioned.NewForConfig(k8sConfig)
	if err != nil {
		return err
	}

	config := &controllers.Config{}
	if err := extension.ReadConfig(configFile, config); err != nil {
		return err
	}

	ext := &controllers.TraefikExtension{
		Traefik: client,
	}

	ext.AddAllowedDomains(config.AllowedDomains)

	fmt.Println("Listening on", config.ListenAddress)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	return extension.Serve[*controllers.Traefik](ctx, config, ext)
}
