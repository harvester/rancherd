package bootstrap

import (
	"github.com/harvester/rancherd/pkg/rancherd"
	"github.com/spf13/cobra"
)

func NewBootstrap() *cobra.Command {
	b := &Bootstrap{}
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Run Rancher and Kubernetes bootstrap",
		RunE:  b.Run,
	}

	cmd.Flags().BoolVarP(&b.Force, "force", "f", false, "Run bootstrap even if already bootstrapped")

	return cmd
}

type Bootstrap struct {
	Force bool
}

func (b *Bootstrap) Run(cmd *cobra.Command, _ []string) error {
	r := rancherd.New(rancherd.Config{
		Force:      b.Force,
		DataDir:    rancherd.DefaultDataDir,
		ConfigPath: rancherd.DefaultConfigFile,
	})
	return r.Run(cmd.Context())
}
