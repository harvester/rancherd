package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/rancher/rancherd/cmd/rancherd/bootstrap"
	"github.com/rancher/rancherd/cmd/rancherd/gettoken"
	"github.com/rancher/rancherd/cmd/rancherd/gettpmhash"
	"github.com/rancher/rancherd/cmd/rancherd/info"
	"github.com/rancher/rancherd/cmd/rancherd/probe"
	"github.com/rancher/rancherd/cmd/rancherd/resetadmin"
	"github.com/rancher/rancherd/cmd/rancherd/retry"
	"github.com/rancher/rancherd/cmd/rancherd/updateclientsecret"
	"github.com/rancher/rancherd/cmd/rancherd/upgrade"
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
		updateclientsecret.NewUpdateClientSecret(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
