package runtime

import (
	"fmt"
	"os"

	"github.com/harvester/rancherd/pkg/config"
	"github.com/harvester/rancherd/pkg/images"
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

func ToUpgradeInstruction(cfg *config.Config, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	instruction, err := ToInstruction(cfg.RuntimeInstallerImage, cfg.SystemDefaultRegistry, k8sVersion)
	if err != nil {
		return nil, err
	}
	return instruction, nil
}

func ToJoinInstruction(imageOverride string, systemDefaultRegistry string, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	runtime := config.GetRuntime(k8sVersion)
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = string(runtime)
	instruction.Env = []string{
		"RESTART_STAMP=" + images.GetInstallerImage(imageOverride, systemDefaultRegistry, k8sVersion),
		"INSTALL_RKE2_TYPE=agent",
	}
	instruction.Image = images.GetInstallerImage(imageOverride, systemDefaultRegistry, k8sVersion)
	instruction.SaveOutput = true
	return instruction, nil
}

func ToStartRKE2Agent(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "start-agent"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", "systemctl", "start", "rke2-agent"}
	instruction.Command = cmd
	return instruction, nil
}

func ToEnableRKE2Agent(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "enable-agent"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", "systemctl", "enable", "rke2-agent"}
	instruction.Command = cmd
	return instruction, nil
}

func ToDisableRKE2Agent(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "disable-agent"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", "systemctl", "disable", "rke2-agent"}
	instruction.Command = cmd
	return instruction, nil
}

func ToStopRKE2Agent(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "stop-agent"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", "systemctl", "stop", "rke2-agent"}
	instruction.Command = cmd
	return instruction, nil
}

func ToEnableRKE2Server(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "enable-server"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", "systemctl", "enable", "rke2-server"}
	instruction.Command = cmd
	return instruction, nil
}

func ToStartRKE2Server(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "start-server"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", "systemctl", "start", "rke2-server"}
	instruction.Command = cmd
	return instruction, nil
}
