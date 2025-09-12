package probe

import (
	"fmt"
	"time"

	"github.com/harvester/rancherd/pkg/probe"
	"github.com/spf13/cobra"
)

func NewProbe() *cobra.Command {
	p := &Probe{}
	cmd := &cobra.Command{
		Use:    "probe",
		Short:  "Run plan probes",
		Hidden: true,
		RunE:   p.Run,
	}

	cmd.Flags().StringVarP(&p.Interval, "interval", "i", "2s", "Polling interval to run probes")
	cmd.Flags().StringVarP(&p.File, "file", "f", "/var/lib/rancher/rancherd/plan/plan.json", "Plan file")

	return cmd
}

type Probe struct {
	Interval string
	File     string
}

func (p *Probe) Run(cmd *cobra.Command, _ []string) error {
	interval, err := time.ParseDuration(p.Interval)
	if err != nil {
		return fmt.Errorf("parsing duration %s: %w", p.Interval, err)
	}

	return probe.RunProbes(cmd.Context(), p.File, interval)
}
