package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/RobertoCCC/postcode-pt-cli/internal/api"
	"github.com/RobertoCCC/postcode-pt-cli/internal/render"
)

var (
	flagAPIURL  string
	flagJSON    bool
	flagNoColor bool
)

var rootCmd = &cobra.Command{
	Use:           "pcpt <postal-code>",
	Short:         "Search Portuguese postal codes from the command line",
	Long:          "pcpt is a CLI for the postcode-pt API. Look up a postal code (CP4-CP3) or browse districts and municipalities.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		cp4, cp3, err := api.NormalizePostalCode(args[0])
		if err != nil {
			return fmt.Errorf("invalid postal code %q: expected CP4-CP3 (e.g. 1100-038)", args[0])
		}
		client := api.NewClient(flagAPIURL)
		entries, err := client.Lookup(cmd.Context(), cp4, cp3)
		if err != nil {
			if errors.Is(err, api.ErrNotFound) {
				return fmt.Errorf("postal code %s-%s not found", cp4, cp3)
			}
			return err
		}
		opts := render.NewOptions(cmd.OutOrStdout())
		opts.JSON = flagJSON
		opts.NoColor = flagNoColor
		return render.PostalCodes(opts, entries)
	},
}

func Execute() {
	if err := rootCmd.ExecuteContext(rootContext()); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "Override the API base URL (defaults to "+api.DefaultBaseURL+")")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Emit raw JSON output")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "Disable ANSI colors")

	rootCmd.AddCommand(districtsCmd)
	rootCmd.AddCommand(districtCmd)
	rootCmd.AddCommand(versionCmd)
}
