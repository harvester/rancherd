package runtime

import (
	"fmt"
	"os"

	"github.com/harvester/rancherd/pkg/kubectl"
	"github.com/harvester/rancherd/pkg/self"
	"github.com/rancher/system-agent/pkg/applyinator"
)

func ToWaitKubernetesInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "wait-harvester-provisioned"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "harvester-system", "rollout", "status", "-w", "deploy/harvester"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}
