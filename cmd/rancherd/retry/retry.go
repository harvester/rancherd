package retry

import (
	"time"

	"github.com/rancher/rancherd/pkg/retry"
	"github.com/spf13/cobra"
)

func NewRetry() *cobra.Command {
	r := &Retry{}
	cmd := &cobra.Command{
		Use:                "retry",
		Short:              "Retry command until it succeeds",
		DisableFlagParsing: true,
		Hidden:             true,
		RunE:               r.Run,
	}

	cmd.Flags().BoolVar(&r.SleepFirst, "sleep-first", false, "Sleep 5 seconds before running command")

	return cmd
}

type Retry struct {
	SleepFirst bool
}

func (p *Retry) Run(cmd *cobra.Command, args []string) error {
	if p.SleepFirst {
		time.Sleep(5 * time.Second)
	}
	return retry.Retry(cmd.Context(), 15*time.Second, args)
}
