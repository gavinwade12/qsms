package cmd

import (
	"bufio"
	"fmt"
	"log"
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
	Short: "send an sms",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		carrier, err := cmd.Flags().GetString("carrier")
		if err != nil {
			log.Fatal(err)
		}
		carrier = strings.ToLower(carrier)

		domain := gateway[carrier]
		if domain == "" {
			log.Fatalf("no domain for carrier: %s", carrier)
		}

		number, err := cmd.Flags().GetString("number")
		if err != nil {
			log.Fatal(err)
		}

		text, err := cmd.Flags().GetString("text")
		if err != nil {
			log.Fatal(err)
		}

		num, err := phonenumbers.Parse(number, "US")
		if err != nil {
			log.Fatal(err)
		}
		number = phonenumbers.Format(num, phonenumbers.E164)[2:]

		from := viper.GetString("email")
		if from == "" {
			fmt.Printf("Email: ")
			r := bufio.NewReader(os.Stdin)
			l, _, err := r.ReadLine()
			if err != nil {
				log.Fatal(err)
			}

			from = string(l)
			if from == "" {
				log.Fatal("an email is required!")
			}

			viper.Set("email", from)
			if err = viper.WriteConfig(); err != nil {
				log.Fatal(err)
			}
		}

		pass := viper.GetString("password")
		if pass == "" {
			fmt.Printf("Password: ")
			r := bufio.NewReader(os.Stdin)
			l, _, err := r.ReadLine()
			if err != nil {
				log.Fatal(err)
			}

			pass = string(l)
			if pass == "" {
				log.Fatal("a password is required!")
			}
			viper.Set("password", pass)
			if err = viper.WriteConfig(); err != nil {
				log.Fatal(err)
			}
		}

		to := fmt.Sprintf("%s@%s", number, gateway[carrier])
		err = smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
			from,
			[]string{to},
			[]byte(fmt.Sprintf("From: %s\r\nTo: %s\r\n\r\n%s", from, to, text)))
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("carrier", "c", "", "carrier e.g. 'Verizon'")
	sendCmd.MarkFlagRequired("carrier")
	sendCmd.Flags().StringP("number", "n", "", "phone number with country code")
	sendCmd.MarkFlagRequired("number")
	sendCmd.Flags().StringP("text", "t", "", "the sms body")
	sendCmd.MarkFlagRequired("text")
}
