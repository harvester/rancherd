package runtime

import (
	"github.com/harvester/rancherd/pkg/config"
	"github.com/harvester/rancherd/pkg/images"
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

func ToUpgradeInstruction(cfg *config.Config, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	instruction, err := ToInstruction(cfg.RuntimeInstallerImage, cfg.SystemDefaultRegistry, k8sVersion)
	if err != nil {
		return nil, err
	}
	return instruction, nil
}
