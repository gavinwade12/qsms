package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
	"github.com/spf13/viper"
)

// TwilioGateway is a Gateway for sending SMS via Twilio's REST API
type TwilioGateway struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// Send implements Gateway.Send
func (g TwilioGateway) Send(recipient, text string) error {
	if g.RequiresConfiguration() {
		return fmt.Errorf("the twilio gateway is not configured")
	}

	num, err := phonenumbers.Parse(recipient, "US")
	if err != nil {
		return err
	}
	recipient = phonenumbers.Format(num, phonenumbers.INTERNATIONAL)

	num, err = phonenumbers.Parse(g.FromNumber, "US")
	if err != nil {
		return err
	}
	from := phonenumbers.Format(num, phonenumbers.INTERNATIONAL)

	msgData := url.Values{}
	msgData.Set("To", recipient)
	msgData.Set("From", from)
	msgData.Set("Body", text)
	msgDataReader := strings.NewReader(msgData.Encode())

	url := "https://api.twilio.com/2010-04-01/Accounts/" + g.AccountSID + "/Messages.json"
	req, err := http.NewRequest(http.MethodPost, url, msgDataReader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(g.AccountSID, g.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: time.Second * 2}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body of request (status: %d): %v", resp.StatusCode, err)
	}
	return fmt.Errorf("got status code: %d\nbody: %s", resp.StatusCode, data)

}

// RequiresConfiguration implements Gateway.RequiresConfiguration
func (g TwilioGateway) RequiresConfiguration() bool {
	return g.AccountSID == "" || g.AuthToken == "" || g.FromNumber == ""
}

// PromptForConfiguration implements Gateway.PromptForConfiguration
func (g *TwilioGateway) PromptForConfiguration() error {
	fmt.Printf("account sid: ")
	r := bufio.NewReader(os.Stdin)
	l, _, err := r.ReadLine()
	if err != nil {
		return err
	}
	g.AccountSID = string(l)
	if g.AccountSID == "" {
		return fmt.Errorf("an account sid is required")
	}
	viper.Set("gateways.twilio.account_sid", g.AccountSID)

	fmt.Printf("auth token: ")
	l, _, err = r.ReadLine()
	if err != nil {
		return err
	}
	g.AuthToken = string(l)
	if g.AuthToken == "" {
		return fmt.Errorf("an auth token is required")
	}
	viper.Set("gateways.twilio.auth_token", g.AuthToken)

	fmt.Printf("number: ")
	l, _, err = r.ReadLine()
	if err != nil {
		return err
	}
	g.FromNumber = string(l)
	if g.FromNumber == "" {
		return fmt.Errorf("a number is required")
	}
	viper.Set("gateways.twilio.from_number", g.FromNumber)

	return viper.WriteConfig()
}
