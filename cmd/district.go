package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RobertoCCC/postcode-pt-cli/internal/api"
	"github.com/RobertoCCC/postcode-pt-cli/internal/render"
)

var districtCmd = &cobra.Command{
	Use:   "district <code>",
	Short: "List municipalities for a given district code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(flagAPIURL)
		municipalities, err := client.ListMunicipalities(cmd.Context(), args[0])
		if err != nil {
			if errors.Is(err, api.ErrNotFound) {
				return fmt.Errorf("district %s not found", args[0])
			}
			return err
		}
		opts := render.NewOptions(cmd.OutOrStdout())
		opts.JSON = flagJSON
		opts.NoColor = flagNoColor
		return render.Municipalities(opts, municipalities)
	},
}
