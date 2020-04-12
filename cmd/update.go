package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the domain/carrier gateway list",
	Long: `To add/edit a carrier, specify the carrier and domain:
qsms update -c verizon -d "vtext.com"

To remove a carrier from the list, use the remove flag:
qsms update -c verizon -r`,
	Run: func(cmd *cobra.Command, args []string) {
		carrier, err := cmd.Flags().GetString("carrier")
		if err != nil {
			exitWithError(err)
		}
		carrier = strings.ToLower(carrier)

		domain, err := cmd.Flags().GetString("domain")
		if err != nil {
			exitWithError(err)
		}

		remove, err := cmd.Flags().GetBool("remove")
		if err != nil {
			exitWithError(err)
		}

		if remove && domain != "" {
			exitWithErrorMessage("The remove and domain flags cannot be present simultaneously.")
		}

		if remove {
			delete(gateway, carrier)
		} else {
			gateway[carrier] = domain
		}

		viper.Set("gateway", gateway)
		if err := viper.WriteConfig(); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("carrier", "c", "", "carrier to update")
	updateCmd.MarkFlagRequired("carrier")
	updateCmd.Flags().StringP("domain", "d", "", "the carrier's domain")
	updateCmd.Flags().BoolP("remove", "r", false, "remove the specificed carrier")

}
