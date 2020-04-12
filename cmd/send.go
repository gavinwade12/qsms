// +build !darwin

package cmd

import (
	"bufio"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/nyaruka/phonenumbers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an sms",
	Long: `Sending an sms requires a carrier, number, and text:
qsms send -c verizon -n 41955512345 -t "Hello World!"`,
	Run: func(cmd *cobra.Command, args []string) {
		carrier, err := cmd.Flags().GetString("carrier")
		if err != nil {
			exitWithError(err)
		}
		carrier = strings.ToLower(carrier)

		domain := gateway[carrier]
		if domain == "" {
			exitWithErrorMessage("no domain for carrier: %s", carrier)
		}

		number, err := cmd.Flags().GetString("number")
		if err != nil {
			exitWithError(err)
		}

		text, err := cmd.Flags().GetString("text")
		if err != nil {
			exitWithError(err)
		}

		num, err := phonenumbers.Parse(number, "US")
		if err != nil {
			exitWithError(err)
		}
		number = phonenumbers.Format(num, phonenumbers.E164)[2:]

		from := viper.GetString("email")
		if from == "" {
			fmt.Printf("Email: ")
			r := bufio.NewReader(os.Stdin)
			l, _, err := r.ReadLine()
			if err != nil {
				exitWithError(err)
			}

			from = string(l)
			if from == "" {
				exitWithErrorMessage("an email is required!")
			}

			viper.Set("email", from)
			if err = viper.WriteConfig(); err != nil {
				exitWithError(err)
			}
		}

		pass := viper.GetString("password")
		if pass == "" {
			fmt.Printf("Password: ")
			r := bufio.NewReader(os.Stdin)
			l, _, err := r.ReadLine()
			if err != nil {
				exitWithError(err)
			}

			pass = string(l)
			if pass == "" {
				exitWithErrorMessage("a password is required!")
			}
			viper.Set("password", pass)
			if err = viper.WriteConfig(); err != nil {
				exitWithError(err)
			}
		}

		to := fmt.Sprintf("%s@%s", number, gateway[carrier])
		err = smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
			from,
			[]string{to},
			[]byte(fmt.Sprintf("From: %s\r\nTo: %s\r\n\r\n%s", from, to, text)))
		if err != nil {
			exitWithErrorMessage("There was an issue sending the message. Is your email and password configured conrrectly? "+
				"You can verify them here: %s.\n\n%v", viper.ConfigFileUsed(), err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("carrier", "c", "", "the recipient carrier")
	sendCmd.MarkFlagRequired("carrier")
	sendCmd.Flags().StringP("number", "n", "", "the recipient phone number")
	sendCmd.MarkFlagRequired("number")
	sendCmd.Flags().StringP("text", "t", "", "the sms body")
	sendCmd.MarkFlagRequired("text")
}
