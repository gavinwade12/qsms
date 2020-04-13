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
	SMTPServer           string
	SMTPServerPort       int
	Sender               string
	SenderPassword       string
	CarrierDomainMapping map[string]string
}

// Send implements Gateway.Send
func (g EmailGateway) Send(recipient, text string) error {
	if g.RequiresConfiguration() {
		return fmt.Errorf("the email gateway is not configured")
	}

	num, err := phonenumbers.Parse(recipient, "US")
	if err != nil {
		return err
	}
	recipient = phonenumbers.Format(num, phonenumbers.E164)[2:]

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

	domain := g.CarrierDomainMapping[carrier]
	if domain == "" {
		return fmt.Errorf("a domain could not be found for the carrier: %s", carrier)
	}

	to := fmt.Sprintf("%s@%s", recipient, domain)
	err = smtp.SendMail(fmt.Sprintf("%s:%d", g.SMTPServer, g.SMTPServerPort),
		smtp.PlainAuth("", g.Sender, g.SenderPassword, g.SMTPServer),
		g.Sender,
		[]string{to},
		[]byte(fmt.Sprintf("From: %s\r\nTo: %s\r\n\r\n%s", g.Sender, to, text)))
	if err != nil {
		return fmt.Errorf("There was an issue sending the message. Is your email gateway configured conrrectly? "+
			"You can verify it here: %s.\n\n%v", viper.ConfigFileUsed(), err)
	}
	return nil
}

// RequiresConfiguration implements Gateway.RequiresConfiguration
func (g EmailGateway) RequiresConfiguration() bool {
	return g.SMTPServer == "" || g.SMTPServerPort <= 0 ||
		g.Sender == "" || g.SenderPassword == "" ||
		len(g.CarrierDomainMapping) == 0
}

// PromptForConfiguration implements Gateway.PromptForConfiguration
func (g *EmailGateway) PromptForConfiguration() error {
	r := bufio.NewReader(os.Stdin)
	if g.Sender == "" {
		fmt.Printf("email: ")
		l, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		g.Sender = string(l)
		if g.Sender == "" {
			return fmt.Errorf("an email is required")
		}
		viper.Set("gateways.email.sender", g.Sender)
	}

	if g.SenderPassword == "" {
		fmt.Printf("password: ")
		l, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		pass := string(l)
		if pass == "" {
			return fmt.Errorf("a password is required")
		}
		g.SenderPassword = base32.StdEncoding.EncodeToString([]byte(pass))
		viper.Set("gateways.email.sender_password", g.SenderPassword)
	}

	if g.SMTPServer == "" {
		fmt.Printf("smtp server: ")
		l, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		g.SMTPServer = string(l)
		if g.SMTPServer == "" {
			return fmt.Errorf("an smtp server is required")
		}
		viper.Set("gateways.email.smtp_server", g.SMTPServer)
	}

	if g.SMTPServerPort <= 0 {
		fmt.Printf("smtp server port: ")
		l, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		port := string(l)
		if port == "" {
			return fmt.Errorf("an smtp server port is required")
		}
		g.SMTPServerPort, err = strconv.Atoi(port)
		if err != nil {
			return err
		}
		viper.Set("gateways.email.smtp_server_port", g.SMTPServerPort)
	}

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	if len(g.CarrierDomainMapping) == 0 {
		return fmt.Errorf(
			"no carriers are configured for the email gateway. please configure them in %s",
			viper.ConfigFileUsed())
	}
	return nil
}
