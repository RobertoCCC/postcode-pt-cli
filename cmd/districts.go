package cmd

import (
	"github.com/spf13/cobra"

	"github.com/RobertoCCC/postcode-pt-cli/internal/api"
	"github.com/RobertoCCC/postcode-pt-cli/internal/render"
)

var districtsCmd = &cobra.Command{
	Use:   "districts",
	Short: "List every Portuguese district",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(flagAPIURL)
		districts, err := client.ListDistricts(cmd.Context())
		if err != nil {
			return err
		}
		opts := render.NewOptions(cmd.OutOrStdout())
		opts.JSON = flagJSON
		opts.NoColor = flagNoColor
		return render.Districts(opts, districts)
	},
}
