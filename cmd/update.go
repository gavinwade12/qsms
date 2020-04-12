package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update the domain/carrier gateway list",
	Long: `By default, an attempt will be made to update the carrier/domain gateway list automatically. The list is pulled from 'https://support.teamunify.com/en/articles/227-email-to-sms-gateway-list'.

If you'd like to add/edit a carrier, you can manually specify the carrier and domain:
qsms update -c "Verizon" -d "vtext.com"

To remove a carrier from the list, use the remove flag:
qsms update -c "Verizon" -r`,
	Run: func(cmd *cobra.Command, args []string) {
		carrier, err := cmd.Flags().GetString("carrier")
		if err != nil {
			log.Fatal(err)
		}
		carrier = strings.ToLower(carrier)

		domain, err := cmd.Flags().GetString("domain")
		if err != nil {
			log.Fatal(err)
		}

		remove, err := cmd.Flags().GetBool("remove")
		if err != nil {
			log.Fatal(err)
		}

		if remove && domain != "" {
			log.Fatal("The remove and domain flags cannot be present simultaneously.")
		}

		if remove {
			delete(gateway, carrier)
		} else {
			gateway[carrier] = domain
		}

		viper.Set("gateway", gateway)
		if err := viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("carrier", "c", "", "carrier e.g. 'Verizon'")
	updateCmd.MarkFlagRequired("carrier")
	updateCmd.Flags().StringP("domain", "d", "", "domain e.g. 'vtext.com'")
	updateCmd.Flags().BoolP("remove", "r", false, "remove the specificed carrier")

}
