package rancher

import (
	"encoding/base64"
	"fmt"

	"github.com/harvester/rancherd/pkg/config"
	"github.com/harvester/rancherd/pkg/images"
	"github.com/harvester/rancherd/pkg/kubectl"
	"github.com/rancher/system-agent/pkg/applyinator"
	"github.com/rancher/wrangler/v3/pkg/data"
	"sigs.k8s.io/yaml"
)

var defaultValues = map[string]interface{}{
	"ingress": map[string]interface{}{
		"enabled": false,
	},
	"features":     "multi-cluster-management=false",
	"antiAffinity": "required",
	"replicas":     -3,
	"tls":          "external",
	"hostPort":     8443,
}

func GetRancherValues(dataDir string) string {
	return fmt.Sprintf("%s/rancher/values.yaml", dataDir)
}

func ToFile(cfg *config.Config, dataDir string) (*applyinator.File, error) {
	values := data.MergeMaps(defaultValues, map[string]interface{}{
		"systemDefaultRegistry": cfg.SystemDefaultRegistry,
	})
	values = data.MergeMaps(values, cfg.RancherValues)

	data, err := yaml.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("marshalling Rancher values.yaml: %w", err)
	}

	return &applyinator.File{
		Content: base64.StdEncoding.EncodeToString(data),
		Path:    GetRancherValues(dataDir),
	}, nil
}

func ToInstruction(imageOverride, systemDefaultRegistry, k8sVersion, rancherVersion, dataDir string) (*applyinator.OneTimeInstruction, error) {
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "rancher"
	instruction.SaveOutput = true
	instruction.Image = images.GetRancherInstallerImage(imageOverride, systemDefaultRegistry, rancherVersion)
	instruction.Env = append(kubectl.Env(k8sVersion), fmt.Sprintf("RANCHER_VALUES=%s", GetRancherValues(dataDir)))
	return instruction, nil
}

func ToUpgradeInstruction(imageOverride, systemDefaultRegistry, k8sVersion, rancherVersion, _ string) (*applyinator.OneTimeInstruction, error) {
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "rancher"
	instruction.SaveOutput = true
	instruction.Image = images.GetRancherInstallerImage(imageOverride, systemDefaultRegistry, rancherVersion)
	instruction.Env = kubectl.Env(k8sVersion)
	return instruction, nil
}
