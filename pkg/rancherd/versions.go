package rancherd

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/harvester/rancherd/pkg/kubectl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func (r *Rancherd) getExistingVersions(ctx context.Context) (k8sVersion, rancherOSVersion string) {
	kubeConfig, err := kubectl.GetKubeconfig("")
	if err != nil {
		return "", ""
	}

	data, err := os.ReadFile(kubeConfig)
	if err != nil {
		return "", ""
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(data)
	if err != nil {
		return "", ""
	}

	k8s, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return "", ""
	}
	return getK8sVersion(ctx, k8s), getRancherOSVersion()
}

func getK8sVersion(ctx context.Context, k8s kubernetes.Interface) string {
	nodes, err := k8s.CoreV1().Nodes().List(ctx, metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/control-plane=true",
	})
	if err != nil || len(nodes.Items) == 0 {
		return ""
	}
	return nodes.Items[0].Status.NodeInfo.KubeletVersion
}

func getRancherOSVersion() string {
	data, err := os.ReadFile("/usr/lib/rancheros-release")
	if err != nil {
		return ""
	}

	scan := bufio.NewScanner(bytes.NewBuffer(data))
	for scan.Scan() {
		if strings.HasPrefix(scan.Text(), "IMAGE=") {
			return strings.TrimSuffix(strings.TrimPrefix(scan.Text(), "IMAGE="), "-"+runtime.GOARCH)
		}
	}
	return ""
}
