package rancher

import (
	"fmt"
	"os"

	"github.com/rancher/system-agent/pkg/applyinator"

	"github.com/harvester/rancherd/pkg/kubectl"
	"github.com/harvester/rancherd/pkg/self"
)

func ToWaitRancherInstruction(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "wait-rancher",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "status", "-w", "deploy/rancher"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

func ToWaitRancherWebhookInstruction(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "wait-rancher-webhook",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "status", "-w", "deploy/rancher-webhook"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

func ToWaitSUCInstruction(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "wait-system-upgrade-controller",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "status", "-w", "deploy/system-upgrade-controller"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

func ToWaitSUCPlanInstruction(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "wait-suc-plan-resolved",
		SaveOutput: true,
		Args: []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "wait",
			"--for=condition=LatestResolved=true", "plans.upgrade.cattle.io", "system-agent-upgrader"},
		Env:     kubectl.Env(k8sVersion),
		Command: cmd,
	}, nil
}

func ToWaitClusterClientSecretInstruction(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "wait-cluster-client-secret-resolved",
		SaveOutput: true,
		Args: []string{"retry", kubectl.Command(k8sVersion), "-n", clusterNamespace, "get",
			"secret", clusterClientSecret},
		Env:     kubectl.Env(k8sVersion),
		Command: cmd,
	}, nil
}

func ToUpdateClientSecretInstruction(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "update-client-secret",
		SaveOutput: true,
		Args:       []string{"update-client-secret"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

func ToScaleDownFleetControllerInstruction(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "scale-down-fleet-controller",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-fleet-system", "scale", "--replicas", "0", "deploy/fleet-controller"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

func ToScaleUpFleetControllerInstruction(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "scale-up-fleet-controller",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-fleet-system", "scale", "--replicas", "1", "deploy/fleet-controller"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

func ToDeleteRancherWebhookValidationConfiguration(k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "delete-rancher-webhook-validation-configuration",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "delete", "validatingwebhookconfiguration", "rancher.cattle.io"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

func ToRestartRancherWebhookInstruction(k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "wait-rancher-webhook",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "-n", "cattle-system", "rollout", "restart", "deploy/rancher-webhook"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

func ToCreateStvAggregationSecret(k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "create-stv-aggregation-secret",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "create", "secret", "generic", "-n", "cattle-system", "stv-aggregation"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}

// Needs to patch status subresource
// k patch cluster.provisioning local -n fleet-local --subresource=status --type=merge --patch '{"status":{"fleetWorkspaceName": "fleet-local"}}'
func PatchLocalProvisioningClusterStatus(_, _, k8sVersion string) (*applyinator.Instruction, error) {
	cmd, err := self.Self()
	if err != nil {
		return nil, fmt.Errorf("resolving location of %s: %w", os.Args[0], err)
	}
	return &applyinator.Instruction{
		Name:       "patch-provisioning-cluster-status",
		SaveOutput: true,
		Args:       []string{"retry", kubectl.Command(k8sVersion), "-n", "fleet-local", "patch", "cluster.provisioning", "local", "--subresource=status", "--type=merge", "--patch", "{\"status\":{\"fleetWorkspaceName\": \"fleet-local\"}}"},
		Env:        kubectl.Env(k8sVersion),
		Command:    cmd,
	}, nil
}
