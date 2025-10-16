package runtime

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/harvester/rancherd/pkg/config"
	"github.com/harvester/rancherd/pkg/images"
	"github.com/harvester/rancherd/pkg/kubectl"
	"github.com/harvester/rancherd/pkg/self"
	"github.com/rancher/system-agent/pkg/applyinator"
)

func ToInstruction(imageOverride string, systemDefaultRegistry string, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	runtime := config.GetRuntime(k8sVersion)
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = string(runtime)
	instruction.Env = []string{
		"RESTART_STAMP=" + images.GetInstallerImage(imageOverride, systemDefaultRegistry, k8sVersion),
	}
	instruction.Image = images.GetInstallerImage(imageOverride, systemDefaultRegistry, k8sVersion)
	instruction.SaveOutput = true
	return instruction, nil
}

func ToUpgradeInstruction(k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	patch, err := json.Marshal(map[string]interface{}{
		"spec": map[string]interface{}{
			"kubernetesVersion": k8sVersion,
		},
	})
	if err != nil {
		return nil, err
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "patch-kubernetes-version"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "--type=merge", "-n", "fleet-local", "patch", "clusters.provisioning.cattle.io", "local", "-p", string(patch)}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}
