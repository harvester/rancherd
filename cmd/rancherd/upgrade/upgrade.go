package upgrade

import (
	"github.com/rancher/rancherd/pkg/rancherd"
	"github.com/spf13/cobra"
)

func NewUpgrade() *cobra.Command {
	u := &Upgrade{}
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade Rancher and Kubernetes",
		RunE:  u.Run,
	}

	cmd.Flags().StringVarP(&u.RancherVersion, "rancher-version", "r", "stable", "Target Rancher version")
	cmd.Flags().StringVarP(&u.RancherOSVersion, "rancher-os-version", "o", "latest", "Target RancherOS version")
	cmd.Flags().StringVarP(&u.KubernetesVersion, "kubernetes-version", "k", "stable", "Target Kubernetes version")
	cmd.Flags().BoolVarP(&u.Force, "force", "f", false, "Run without prompting for confirmation")

	return cmd
}

type Upgrade struct {
	RancherVersion    string
	RancherOSVersion  string
	KubernetesVersion string
	Force             bool
}

func (b *Upgrade) Run(cmd *cobra.Command, _ []string) error {
	r := rancherd.New(rancherd.Config{
		Force:      b.Force,
		DataDir:    rancherd.DefaultDataDir,
		ConfigPath: rancherd.DefaultConfigFile,
	})
	return r.Upgrade(cmd.Context(), rancherd.UpgradeConfig{
		RancherVersion:    b.RancherVersion,
		KubernetesVersion: b.KubernetesVersion,
		RancherOSVersion:  b.RancherOSVersion,
	})
}
