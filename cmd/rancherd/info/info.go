package info

import (
	"github.com/harvester/rancherd/pkg/rancherd"
	"github.com/spf13/cobra"
)

func NewInfo() *cobra.Command {
	i := &Info{}
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Print installation versions",
		RunE:  i.Run,
	}

	return cmd
}

type Info struct {
}

func (b *Info) Run(cmd *cobra.Command, _ []string) error {
	r := rancherd.New(rancherd.Config{
		DataDir:    rancherd.DefaultDataDir,
		ConfigPath: rancherd.DefaultConfigFile,
	})
	return r.Info(cmd.Context())
}
