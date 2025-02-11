package resources

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	v1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/system-agent/pkg/applyinator"
	"github.com/rancher/wrangler/pkg/data/convert"
	"github.com/rancher/wrangler/pkg/randomtoken"
	"github.com/rancher/wrangler/pkg/yaml"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rancher/rancherd/pkg/config"
	"github.com/rancher/rancherd/pkg/images"
	"github.com/rancher/rancherd/pkg/kubectl"
	"github.com/rancher/rancherd/pkg/self"
	"github.com/rancher/rancherd/pkg/versions"
)

const (
	localRKEStateSecretType = "rke.cattle.io/cluster-state"
)

func writeCattleID(id string) error {
	if err := os.MkdirAll("/etc/rancher", 0755); err != nil {
		return fmt.Errorf("mkdir /etc/rancher: %w", err)
	}
	if err := os.MkdirAll("/etc/rancher/agent", 0700); err != nil {
		return fmt.Errorf("mkdir /etc/rancher/agent: %w", err)
	}
	return ioutil.WriteFile("/etc/rancher/agent/cattle-id", []byte(id), 0400)
}

func getCattleID() (string, error) {
	data, err := ioutil.ReadFile("/etc/rancher/agent/cattle-id")
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	id := strings.TrimSpace(string(data))
	if id == "" {
		id, err = randomtoken.Generate()
		if err != nil {
			return "", err
		}
		return id, writeCattleID(id)
	}
	return id, nil
}

func ToBootstrapFile(config *config.Config, path string) (*applyinator.File, error) {
	nodeName := config.NodeName
	if nodeName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("looking up hostname: %w", err)
		}
		nodeName = strings.Split(hostname, ".")[0]
	}

	k8sVersion, err := versions.K8sVersion(config.KubernetesVersion)
	if err != nil {
		return nil, err
	}

	token := config.Token
	if token == "" {
		token, err = randomtoken.Generate()
		if err != nil {
			return nil, err
		}
	}

	resources := config.Resources
	return ToFile(append(resources, v1.GenericMap{
		Data: map[string]interface{}{
			"kind":       "Node",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name": nodeName,
				"labels": map[string]interface{}{
					"node-role.kubernetes.io/etcd": "true",
				},
			},
		},
	}, v1.GenericMap{
		Data: map[string]interface{}{
			"kind":       "Namespace",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name": "fleet-local",
			},
		},
	}, v1.GenericMap{
		Data: map[string]interface{}{
			"kind":       "Cluster",
			"apiVersion": "provisioning.cattle.io/v1",
			"metadata": map[string]interface{}{
				"name":      "local",
				"namespace": "fleet-local",
				"labels": map[string]interface{}{
					"provisioning.cattle.io/management-cluster-name": "local",
				},
			},
			"spec": map[string]interface{}{
				"kubernetesVersion": k8sVersion,
				// Rancher needs a non-null rkeConfig to apply system-upgrade-controller managed chart.
				"rkeConfig": map[string]interface{}{},
			},
		},
	}, v1.GenericMap{
		Data: map[string]interface{}{
			"kind":       "Secret",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "local-rke-state",
				"namespace": "fleet-local",
			},
			"type": localRKEStateSecretType,
			"data": map[string]interface{}{
				"serverToken": []byte(token),
				"agentToken":  []byte(token),
			},
		},
	}, v1.GenericMap{
		Data: map[string]interface{}{
			"kind":       "ClusterRegistrationToken",
			"apiVersion": "management.cattle.io/v3",
			"metadata": map[string]interface{}{
				"name":      "default-token",
				"namespace": "local",
			},
			"spec": map[string]interface{}{
				"clusterName": "local",
			},
			"status": map[string]interface{}{
				"token": token,
			},
		},
	}, v1.GenericMap{
		Data: map[string]interface{}{
			"apiVersion": "catalog.cattle.io/v1",
			"kind":       "ClusterRepo",
			"metadata": map[string]interface{}{
				"name": "rancher-stable",
			},
			"spec": map[string]interface{}{
				"url": "https://releases.rancher.com/server-charts/stable",
			},
		},
	}), path)
}

func ToHarvesterClusterRepoFile(path string) (*applyinator.File, error) {
	file := "/usr/share/rancher/rancherd/config.yaml.d/91-harvester-bootstrap-repo.yaml"
	bytes, err := os.ReadFile(file)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	logrus.Infof("Loading config file [%s]", file)
	values := map[string]interface{}{}
	if err := yaml.Unmarshal(bytes, &values); err != nil {
		return nil, err
	}

	result := config.Config{}
	convert.ToObj(values, &result)

	resources := []v1.GenericMap{}
	for _, resource := range result.Resources {
		if resource.Data["kind"] == "Deployment" || resource.Data["kind"] == "Service" {
			resources = append(resources, resource)
		}
	}
	return ToFile(resources, path)
}

func ToFile(resources []v1.GenericMap, path string) (*applyinator.File, error) {
	if len(resources) == 0 {
		return nil, nil
	}

	var objs []runtime.Object
	for _, resource := range resources {
		objs = append(objs, &unstructured.Unstructured{
			Object: resource.Data,
		})
	}

	data, err := yaml.ToBytes(objs)
	if err != nil {
		return nil, err
	}

	return &applyinator.File{
		Content: base64.StdEncoding.EncodeToString(data),
		Path:    path,
	}, nil
}

func GetBootstrapManifests(dataDir string) string {
	return fmt.Sprintf("%s/bootstrapmanifests/rancherd.yaml", dataDir)
}

func ToInstruction(imageOverride, systemDefaultRegistry, k8sVersion, dataDir string) (*applyinator.Instruction, error) {
	bootstrap := GetBootstrapManifests(dataDir)
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "bootstrap",
		SaveOutput: true,
		Image:      images.GetInstallerImage(imageOverride, systemDefaultRegistry, k8sVersion),
		Args:       []string{"retry", kubectl.Command(k8sVersion), "apply", "--validate=false", "-f", bootstrap},
		Command:    cmd,
		Env:        kubectl.Env(k8sVersion),
	}, nil
}

func GetHarvesterClusterRepoManifests(dataDir string) string {
	return fmt.Sprintf("%s/bootstrapmanifests/harvester-cluster-repo.yaml", dataDir)
}

func ToHarvesterClusterRepoInstruction(imageOverride, systemDefaultRegistry, k8sVersion, dataDir string) (*applyinator.Instruction, error) {
	bootstrap := GetHarvesterClusterRepoManifests(dataDir)
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "harvester-cluster-repo",
		SaveOutput: true,
		Image:      images.GetInstallerImage(imageOverride, systemDefaultRegistry, k8sVersion),
		Args:       []string{"retry", kubectl.Command(k8sVersion), "apply", "--validate=false", "-f", bootstrap},
		Command:    cmd,
		Env:        kubectl.Env(k8sVersion),
	}, nil
}

func ToWaitHarvesterClusterRepoInstruction(k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "wait-harvester-cluster-repo",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "status", "-w", "deploy/harvester-cluster-repo"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}
