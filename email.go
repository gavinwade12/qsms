package main

import (
	"bufio"
	"encoding/base32"
	"fmt"
	"net/smtp"
	"os"
	"strconv"

	"github.com/nyaruka/phonenumbers"
	"github.com/spf13/viper"
)

// EmailGateway is a Gateway for sending SMS via email
type EmailGateway struct {
	SMTPServer     string
	Sender         string
	SenderPassword string
}

// Send implements Gateway.Send
func (g EmailGateway) Send(recipient, text string) error {
	num, err := phonenumbers.Parse(recipient, "US")
	if err != nil {
		return err
	}
	recipient = phonenumbers.Format(num, phonenumbers.E164)[2:]

	from := viper.GetString("gateways.email.sender")
	pass := viper.GetString("gateways.email.sender_password")
	smtpServer := viper.GetString("gateways.email.smtp_server")
	smtpServerPort := viper.GetInt("gateways.email.smtp_server_port")
	if smtpServerPort == 0 {
		smtpServerPort = 587
	}
	if from == "" || pass == "" || smtpServer == "" {
		fmt.Println("the email gateway has not been completely configured")
		fmt.Printf("email: ")
		r := bufio.NewReader(os.Stdin)
		l, _, err := r.ReadLine()
		if err != nil {
			return err
		}

		from = string(l)
		if from == "" {
			return fmt.Errorf("an email is required")
		}

		viper.Set("gateways.email.sender", from)

		fmt.Printf("password: ")
		l, _, err = r.ReadLine()
		if err != nil {
			return err
		}

		pass = string(l)
		if pass == "" {
			return fmt.Errorf("a password is required")
		}
		pass = base32.StdEncoding.EncodeToString([]byte(pass))
		viper.Set("gateways.email.sender_password", pass)
		if err = viper.WriteConfig(); err != nil {
			return err
		}

		fmt.Printf("smtp server: ")
		l, _, err = r.ReadLine()
		if err != nil {
			return err
		}

		smtpServer = string(l)
		if smtpServer == "" {
			return fmt.Errorf("an smtp server is required")
		}
		viper.Set("gateways.email.smtp_server", smtpServer)

		fmt.Printf("smtp server port: ")
		l, _, err = r.ReadLine()
		if err != nil {
			return err
		}

		port := string(l)
		if port != "" {
			smtpServerPort, err = strconv.Atoi(port)
			if err != nil {
				return err
			}

			viper.Set("gateways.email.smtp_server_port", smtpServerPort)
		}

		if err = viper.WriteConfig(); err != nil {
			return err
		}
	}

	fmt.Printf("recipient carrier: ")
	r := bufio.NewReader(os.Stdin)
	l, _, err := r.ReadLine()
	if err != nil {
		return err
	}

	carrier := string(l)
	if carrier == "" {
		return fmt.Errorf("a carrier is required")
	}

	mapping := viper.GetStringMapString("gateways.email.mapping")
	domain := mapping[carrier]
	if domain == "" {
		return fmt.Errorf("a domain could not be found for the carrier: %s", carrier)
	}

	to := fmt.Sprintf("%s@%s", recipient, domain)
	p, err := base32.StdEncoding.DecodeString(pass)
	if err != nil {
		return err
	}
	err = smtp.SendMail(fmt.Sprintf("%s:587", smtpServer),
		smtp.PlainAuth("", from, string(p), smtpServer),
		from,
		[]string{to},
		[]byte(fmt.Sprintf("From: %s\r\nTo: %s\r\n\r\n%s", from, to, text)))
	if err != nil {
		return fmt.Errorf("There was an issue sending the message. Is your email gateway configured conrrectly? "+
			"You can verify it here: %s.\n\n%v", viper.ConfigFileUsed(), err)
	}
	return nil
}
