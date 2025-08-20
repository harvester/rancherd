package gettoken

import (
	"fmt"
	"os"

	"github.com/rancher/rancherd/pkg/token"
	"github.com/spf13/cobra"
)

func NewGetToken() *cobra.Command {
	g := &GetToken{}
	cmd := &cobra.Command{
		Use:   "get-token",
		Short: "Print token to join nodes to the cluster",
		RunE:  g.Run,
	}

	cmd.Flags().StringVar(&g.Kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "Kubeconfig file")

	return cmd
}

type GetToken struct {
	Kubeconfig string
}

func (p *GetToken) Run(cmd *cobra.Command, _ []string) error {
	str, err := token.GetToken(cmd.Context(), p.Kubeconfig)
	if err != nil {
		return err
	}
	fmt.Println(str)
	return nil
}
