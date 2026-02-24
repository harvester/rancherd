package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/harvester/rancherd/cmd/rancherd/bootstrap"
	"github.com/harvester/rancherd/cmd/rancherd/gettoken"
	"github.com/harvester/rancherd/cmd/rancherd/gettpmhash"
	"github.com/harvester/rancherd/cmd/rancherd/info"
	"github.com/harvester/rancherd/cmd/rancherd/probe"
	"github.com/harvester/rancherd/cmd/rancherd/resetadmin"
	"github.com/harvester/rancherd/cmd/rancherd/retry"
	"github.com/harvester/rancherd/cmd/rancherd/upgrade"
)

func main() {
	root := &cobra.Command{
		Use:  "rancherd",
		Long: "Bootstrap Rancher and k3s/rke2 on a node",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		bootstrap.NewBootstrap(),
		gettoken.NewGetToken(),
		resetadmin.NewResetAdmin(),
		probe.NewProbe(),
		retry.NewRetry(),
		upgrade.NewUpgrade(),
		info.NewInfo(),
		gettpmhash.NewGetTPMHash(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
