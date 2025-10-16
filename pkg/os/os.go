package os

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/harvester/rancherd/pkg/kubectl"
	"github.com/harvester/rancherd/pkg/self"
	"github.com/rancher/system-agent/pkg/applyinator"
)

func ToUpgradeInstruction(k8sVersion, rancherOSVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	patch, err := json.Marshal(map[string]interface{}{
		"spec": map[string]interface{}{
			"osImage": rancherOSVersion,
		},
	})
	if err != nil {
		return nil, err
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "patch-rancher-os-version"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "--type=merge", "-n", "fleet-local", "patch", "managedosimages.rancheros.cattle.io", "default-os-image", "-p", string(patch)}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}
