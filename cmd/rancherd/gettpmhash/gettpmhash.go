package gettpmhash

import (
	"fmt"

	"github.com/harvester/rancherd/pkg/tpm"
	"github.com/spf13/cobra"
)

func NewGetTPMHash() *cobra.Command {
	g := &GetTPMHash{}
	cmd := &cobra.Command{
		Use:   "get-tpm-hash",
		Short: "Print TPM hash to identify this machine",
		RunE:  g.Run,
	}

	return cmd
}

type GetTPMHash struct {
}

func (p *GetTPMHash) Run(*cobra.Command, []string) error {
	str, err := tpm.GetPubHash()
	if err != nil {
		return err
	}
	fmt.Println(str)
	return nil
}
