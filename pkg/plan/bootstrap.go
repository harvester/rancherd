package plan

import (
	"context"
	"fmt"

	"github.com/rancher/system-agent/pkg/applyinator"

	"github.com/harvester/rancherd/pkg/config"
	"github.com/harvester/rancherd/pkg/discovery"
	"github.com/harvester/rancherd/pkg/kubectl"
	"github.com/harvester/rancherd/pkg/probe"
	"github.com/harvester/rancherd/pkg/registry"
	"github.com/harvester/rancherd/pkg/resources"
	"github.com/harvester/rancherd/pkg/runtime"
	"github.com/harvester/rancherd/pkg/versions"
)

type plan applyinator.Plan

func toInitPlan(config *config.Config, dataDir string) (*applyinator.Plan, error) {
	if err := assignTokenIfUnset(config); err != nil {
		return nil, err
	}

	plan := plan{}
	if err := plan.addFiles(config, dataDir); err != nil {
		return nil, err
	}

	if err := plan.addInstructions(config, dataDir); err != nil {
		return nil, err
	}

	if err := plan.addProbes(config); err != nil {
		return nil, err
	}

	return (*applyinator.Plan)(&plan), nil
}

func toJoinPlan(cfg *config.Config, dataDir string) (*applyinator.Plan, error) {
	k8sVersion, err := versions.K8sVersion(cfg.KubernetesVersion)
	if err != nil {
		return nil, err
	}
	if cfg.Server == "" {
		return nil, fmt.Errorf("server is required in config for all roles besides cluster-init")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("token is required in config for all roles besides cluster-init")
	}

	plan := plan{}

	runtimeName := config.GetRuntime(k8sVersion)

	// config.yaml
	if err := plan.addFile(runtime.ToJoinFile(&cfg.RuntimeConfig, runtimeName, cfg.Server)); err != nil {
		return nil, err
	}

	// registries.yaml
	if err := plan.addFile(registry.ToFile(cfg.Registries, runtimeName)); err != nil {
		return nil, err
	}

	// install rke2
	if err := plan.addInstruction(runtime.ToJoinInstruction(cfg.RuntimeInstallerImage, cfg.SystemDefaultRegistry, k8sVersion)); err != nil {
		return nil, err
	}

	// enable rke2 agent systemd unit
	if err := plan.addInstruction(runtime.ToEnableRKE2Agent("", "", k8sVersion)); err != nil {
		return nil, err
	}

	// start rke2 agent
	if err := plan.addInstruction(runtime.ToStartRKE2Agent("", "", k8sVersion)); err != nil {
		return nil, err
	}

	if err := plan.addInstruction(probe.ToInstruction()); err != nil {
		return nil, err
	}
	if err := plan.addProbesForJoin(cfg); err != nil {
		return nil, err
	}

	return (*applyinator.Plan)(&plan), nil
}

func ToPlan(ctx context.Context, config *config.Config, dataDir string) (*applyinator.Plan, error) {
	newCfg := *config
	if err := discovery.DiscoverServerAndRole(ctx, &newCfg); err != nil {
		return nil, err
	}
	if newCfg.Role == "cluster-init" {
		return toInitPlan(&newCfg, dataDir)
	}
	return toJoinPlan(&newCfg, dataDir)
}

func (p *plan) addInstructions(cfg *config.Config, dataDir string) error {
	k8sVersion, err := versions.K8sVersion(cfg.KubernetesVersion)
	if err != nil {
		return err
	}

	if err := p.addInstruction(runtime.ToInstruction(cfg.RuntimeInstallerImage, cfg.SystemDefaultRegistry, k8sVersion)); err != nil {
		return err
	}

	if err := p.addInstruction(probe.ToInstruction()); err != nil {
		return err
	}

	if cfg.SkipHarvesterBootstrap {
		return nil
	}

	// If clusterrepo check fails, it waits 5 minutes and retries.
	// Install harvester-cluster-repo deployment before clusterrepo,
	// so we can avoid the 5 minutes waiting time.
	if err := p.addInstruction(resources.ToHarvesterClusterRepoInstruction(cfg.RancherInstallerImage, cfg.SystemDefaultRegistry, k8sVersion, dataDir)); err != nil {
		return err
	}

	if err := p.addInstruction(resources.ToWaitHarvesterClusterRepoInstruction(k8sVersion)); err != nil {
		return err
	}

	if err := p.addInstruction(resources.ToInstruction(cfg.RancherInstallerImage, cfg.SystemDefaultRegistry, k8sVersion, dataDir)); err != nil {
		return err
	}

	p.addPrePostInstructions(cfg, k8sVersion)
	return nil
}

func (p *plan) addPrePostInstructions(cfg *config.Config, k8sVersion string) {
	var instructions []applyinator.OneTimeInstruction

	for _, inst := range cfg.PreInstructions {
		if k8sVersion != "" {
			inst.Env = append(inst.Env, kubectl.Env(k8sVersion)...)
		}
		instructions = append(instructions, inst)
	}

	instructions = append(instructions, p.OneTimeInstructions...)

	for _, inst := range cfg.PostInstructions {
		inst.Env = append(inst.Env, kubectl.Env(k8sVersion)...)
		instructions = append(instructions, inst)
	}

	p.OneTimeInstructions = instructions
}

func (p *plan) addInstruction(instruction *applyinator.OneTimeInstruction, err error) error {
	if err != nil || instruction == nil {
		return err
	}

	p.OneTimeInstructions = append(p.OneTimeInstructions, *instruction)
	return nil
}

func (p *plan) addFiles(cfg *config.Config, dataDir string) error {
	k8sVersions, err := versions.K8sVersion(cfg.KubernetesVersion)
	if err != nil {
		return err
	}
	runtimeName := config.GetRuntime(k8sVersions)

	// config.yaml
	if err := p.addFile(runtime.ToFile(&cfg.RuntimeConfig, runtimeName, true)); err != nil {
		return err
	}

	// bootstrap config.yaml
	if err := p.addFile(runtime.ToBootstrapFile(runtimeName)); err != nil {
		return err
	}

	// registries.yaml
	if err := p.addFile(registry.ToFile(cfg.Registries, runtimeName)); err != nil {
		return err
	}

	if !cfg.SkipHarvesterBootstrap {
		// harvester-cluster-repo manifests
		if err := p.addFile(resources.ToHarvesterClusterRepoFile(resources.GetHarvesterClusterRepoManifests(dataDir))); err != nil {
			return err
		}
	}

	// bootstrap manifests
	return p.addFile(resources.ToBootstrapFile(cfg, resources.GetBootstrapManifests(dataDir)))

}

func (p *plan) addFile(file *applyinator.File, err error) error {
	if err != nil || file == nil {
		return err
	}
	p.Files = append(p.Files, *file)
	return nil
}

func (p *plan) addProbesForJoin(cfg *config.Config) error {
	p.Probes = probe.ProbesForJoin(&cfg.RuntimeConfig)
	return nil
}

func (p *plan) addProbes(cfg *config.Config) error {
	k8sVersion, err := versions.K8sVersion(cfg.KubernetesVersion)
	if err != nil {
		return err
	}
	p.Probes = probe.AllProbes(config.GetRuntime(k8sVersion))
	return nil
}

func ToPromotePlan(cfg *config.Config) (*applyinator.Plan, error) {

	k8sVersion, err := versions.K8sVersion(cfg.KubernetesVersion)
	if err != nil {
		return nil, err
	}
	if cfg.Server == "" {
		return nil, fmt.Errorf("server is required in config for all roles besides cluster-init")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("token is required in config for all roles besides cluster-init")
	}

	plan := plan{}

	runtimeName := config.GetRuntime(k8sVersion)

	// rke2 allows promotion to achieve this we need to do the following
	// 1. generate new /etc/rancher/rke2/config.yaml.d/40-rancherd.yaml which contains role updated to server
	// 2. ensure serverURL is set correctly in 40-rancherd.yaml
	// 3. stop rke2-agent systemd unit
	// 4. disable rke2-agent systemd unit
	// 5. enable rke2-server systemd unit
	// 6. start rke2-server systemd unit

	// generate new 40-rancherd.yaml
	if err := plan.addFile(runtime.ToPromoteFile(&cfg.RuntimeConfig, runtimeName, cfg.Server)); err != nil {
		return nil, err
	}

	// switch node role
	if err := plan.addInstruction(runtime.ToStopRKE2Agent("", "", k8sVersion)); err != nil {
		return nil, err
	}
	if err := plan.addInstruction(runtime.ToDisableRKE2Agent("", "", k8sVersion)); err != nil {
		return nil, err
	}
	if err := plan.addInstruction(runtime.ToEnableRKE2Server("", "", k8sVersion)); err != nil {
		return nil, err
	}
	if err := plan.addInstruction(runtime.ToStartRKE2Server("", "", k8sVersion)); err != nil {
		return nil, err
	}

	return (*applyinator.Plan)(&plan), nil
}
