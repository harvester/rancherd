package plan

import (
	"github.com/harvester/rancherd/pkg/config"
	"github.com/harvester/rancherd/pkg/os"
	"github.com/harvester/rancherd/pkg/runtime"
	"github.com/rancher/system-agent/pkg/applyinator"
)

func Upgrade(cfg *config.Config, k8sVersion, rancherOSVersion, dataDir string) (*applyinator.Plan, error) {
	p := plan{}

	if k8sVersion != "" {
		if err := p.addInstruction(runtime.ToUpgradeInstruction(cfg, k8sVersion)); err != nil {
			return nil, err
		}
	}

	if rancherOSVersion != "" {
		if err := p.addInstruction(os.ToUpgradeInstruction(k8sVersion, rancherOSVersion)); err != nil {
			return nil, err
		}
	}
	return (*applyinator.Plan)(&p), nil
}
