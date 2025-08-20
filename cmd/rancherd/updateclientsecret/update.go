package updateclientsecret

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/rancher/rancherd/pkg/rancher"
)

func NewUpdateClientSecret() *cobra.Command {
	u := &UpdateClientSecret{}
	cmd := &cobra.Command{
		Use:   "update-client-secret",
		Short: "Update cluster client secret to have API Server URL and CA Certs configured",
		RunE:  u.Run,
	}

	cmd.Flags().StringVar(&u.Kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "Kubeconfig file")

	return cmd
}

type UpdateClientSecret struct {
	Kubeconfig string
}

func (s *UpdateClientSecret) Run(cmd *cobra.Command, _ []string) error {
	return rancher.UpdateClientSecret(cmd.Context(), &rancher.Options{Kubeconfig: s.Kubeconfig})
}
