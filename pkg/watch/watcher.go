package watch

import (
	"fmt"

	"os"

	"github.com/harvester/rancherd/pkg/controller"
	"github.com/rancher/wrangler/v3/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/v3/pkg/signals"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	KubeletKubeConfig = "/var/lib/rancher/rke2/agent/kubelet.kubeconfig"
)

func Setup() error {

	val, ok := os.LookupEnv("RANCHERD_DEBUG")
	if ok && val == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = KubeletKubeConfig
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	k8sclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}
	ctx := signals.SetupSignalContext()

	coreFactory, err := core.NewFactoryFromConfig(config)
	if err != nil {
		return fmt.Errorf("error generating core factory: %w", err)
	}

	if err := controller.Register(ctx, coreFactory.Core().V1().Node(), k8sclient.CoreV1().Nodes()); err != nil {
		return fmt.Errorf("error registering controller: %w", err)
	}

	if err := coreFactory.Start(ctx, 1); err != nil {
		return fmt.Errorf("error starting core factory: %w", err)
	}
	<-ctx.Done()
	return nil
}
