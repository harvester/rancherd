package rancher

import (
	"fmt"
	"os"

	"github.com/rancher/system-agent/pkg/applyinator"

	"github.com/harvester/rancherd/pkg/kubectl"
	"github.com/harvester/rancherd/pkg/self"
)

func ToWaitRancherInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "wait-rancher"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "status", "-w", "deploy/rancher"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToWaitRancherWebhookInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "wait-rancher-webhook"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "status", "-w", "deploy/rancher-webhook"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToWaitSUCInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "wait-system-upgrade-controller"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "status", "-w", "deploy/system-upgrade-controller"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToWaitSUCPlanInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "wait-suc-plan-resolved"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "wait",
		"--for=condition=LatestResolved=true", "plans.upgrade.cattle.io", "system-agent-upgrader"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToWaitClusterClientSecretInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "wait-cluster-client-secret-resolved"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", clusterNamespace, "get",
		"secret", clusterClientSecret}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToUpdateClientSecretInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "update-client-secret"
	instruction.SaveOutput = true
	instruction.Args = []string{"update-client-secret"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToScaleDownFleetControllerInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "scale-down-fleet-controller"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-fleet-system", "scale", "--replicas", "0", "deploy/fleet-controller"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToScaleUpFleetControllerInstruction(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "scale-up-fleet-controller"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-fleet-system", "scale", "--replicas", "1", "deploy/fleet-controller"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToDeleteRancherWebhookValidationConfiguration(k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "delete-rancher-webhook-validation-configuration"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "delete", "validatingwebhookconfiguration", "rancher.cattle.io"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToRestartRancherWebhookInstruction(k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "restart-rancher-webhook"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "restart", "deploy/rancher-webhook"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

func ToCreateStvAggregationSecret(k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "create-stv-aggregation-secret"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "create", "secret", "generic", "-n", "cattle-system", "stv-aggregation"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}

// Needs to patch status subresource
// k patch cluster.provisioning local -n fleet-local --subresource=status --type=merge --patch '{"status":{"fleetWorkspaceName": "fleet-local"}}'
func PatchLocalProvisioningClusterStatus(_, _, k8sVersion string) (*applyinator.OneTimeInstruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	instruction := &applyinator.OneTimeInstruction{}
	instruction.Name = "patch-provisioning-cluster-status"
	instruction.SaveOutput = true
	instruction.Args = []string{"retry", kubectl.Command(k8sVersion), "-n", "fleet-local", "patch", "cluster.provisioning", "local", "--subresource=status", "--type=merge", "--patch", "{\"status\":{\"fleetWorkspaceName\": \"fleet-local\"}}"}
	instruction.Env = kubectl.Env(k8sVersion)
	instruction.Command = cmd
	return instruction, nil
}
