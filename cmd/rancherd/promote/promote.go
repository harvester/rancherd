package promote

import (
	"github.com/harvester/rancherd/pkg/rancherd"
	"github.com/spf13/cobra"
)

func NewPromote() *cobra.Command {
	p := &Promote{}
	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote node",
		RunE:  p.Run,
	}

	cmd.Flags().BoolVarP(&p.Force, "force", "f", false, "Run bootstrap even if already bootstrapped")

	return cmd
}

type Promote struct {
	Force bool
}

func (p *Promote) Run(cmd *cobra.Command, _ []string) error {
	r := rancherd.New(rancherd.Config{
		Force:      p.Force,
		DataDir:    rancherd.DefaultDataDir,
		ConfigPath: rancherd.DefaultConfigFile,
	})
	return r.RunPromote(cmd.Context())
}
