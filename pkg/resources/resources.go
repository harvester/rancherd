package resources

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	v1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/system-agent/pkg/applyinator"
	"github.com/rancher/wrangler/v3/pkg/data/convert"
	"github.com/rancher/wrangler/v3/pkg/randomtoken"
	"github.com/rancher/wrangler/v3/pkg/yaml"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/harvester/rancherd/pkg/config"
	"github.com/harvester/rancherd/pkg/images"
	"github.com/harvester/rancherd/pkg/kubectl"
	"github.com/harvester/rancherd/pkg/self"
)

func ToBootstrapFile(config *config.Config, path string) (*applyinator.File, error) {
	nodeName := config.NodeName
	if nodeName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("looking up hostname: %w", err)
		}
		nodeName = strings.Split(hostname, ".")[0]
	}

	token := config.Token
	if token == "" {
		token, err := randomtoken.Generate()
		if err != nil {
			return nil, err
		}
		config.Token = token
	}

	resources := config.Resources
	resources = append(resources, v1.GenericMap{
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
		}})

	return ToFile(resources, path)
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
	if err := convert.ToObj(values, &result); err != nil {
		return nil, err
	}

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

func ToInstruction(imageOverride, systemDefaultRegistry, k8sVersion, dataDir string) (*applyinator.OneTimeInstruction, error) {
	bootstrap := GetBootstrapManifests(dataDir)
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "bootstrap"
	instruction.SaveOutput = true
	instruction.Image = images.GetInstallerImage(imageOverride, systemDefaultRegistry, k8sVersion)
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "apply", "--validate=false", "-f", bootstrap}
	instruction.Command = cmd
	instruction.Env = kubectl.Env(k8sVersion)
	return instruction, nil
}

func GetHarvesterClusterRepoManifests(dataDir string) string {
	return fmt.Sprintf("%s/bootstrapmanifests/harvester-cluster-repo.yaml", dataDir)
}

func ToHarvesterClusterRepoInstruction(imageOverride, systemDefaultRegistry, k8sVersion, dataDir string) (*applyinator.OneTimeInstruction, error) {
	bootstrap := GetHarvesterClusterRepoManifests(dataDir)
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "harvester-cluster-repo"
	instruction.SaveOutput = true
	instruction.Image = images.GetInstallerImage(imageOverride, systemDefaultRegistry, k8sVersion)
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "apply", "--validate=false", "-f", bootstrap}
	instruction.Command = cmd
	instruction.Env = kubectl.Env(k8sVersion)
	return instruction, nil
}

func ToWaitHarvesterClusterRepoInstruction(k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "wait-harvester-cluster-repo"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "status", "-w", "deploy/harvester-cluster-repo"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}
