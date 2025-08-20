package resetadmin

import (
	"os"

	"github.com/rancher/rancherd/pkg/auth"
	"github.com/spf13/cobra"
)

func NewResetAdmin() *cobra.Command {
	r := &ResetAdmin{}
	cmd := &cobra.Command{
		Use:   "reset-admin",
		Short: "Bootstrap and reset admin password",
		RunE:  r.Run,
	}

	cmd.Flags().StringVar(&r.Password, "password", os.Getenv("PASSWORD"), "Password for Rancher login")
	cmd.Flags().StringVar(&r.PasswordFile, "password-file", os.Getenv("PASSWORD_FILE"), "Password for Rancher login, from file")
	cmd.Flags().StringVar(&r.Kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "Kubeconfig file")

	return cmd
}

type ResetAdmin struct {
	Password     string
	PasswordFile string
	Kubeconfig   string
}

func (p *ResetAdmin) Run(cmd *cobra.Command, _ []string) error {
	return auth.ResetAdmin(cmd.Context(), &auth.Options{
		Password:     p.Password,
		PasswordFile: p.PasswordFile,
		Kubeconfig:   p.Kubeconfig,
	})
}
