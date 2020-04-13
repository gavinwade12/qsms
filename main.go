package main

import (
	"bufio"
	"encoding/base32"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var gateways = make(map[string]Gateway)

func init() {
	initConfig()

	senderPass, _ := base32.StdEncoding.DecodeString(viper.GetString("gateways.email.sender_password"))
	gateways["email"] = &EmailGateway{
		SMTPServer:           viper.GetString("gateways.email.smtp_server"),
		SMTPServerPort:       viper.GetInt("gateways.email.smtp_server_port"),
		Sender:               viper.GetString("gateways.email.sender"),
		SenderPassword:       string(senderPass),
		CarrierDomainMapping: viper.GetStringMapString("gateways.email.mapping"),
	}
	gateways["twilio"] = &TwilioGateway{
		AccountSID: viper.GetString("gateways.twilio.account_sid"),
		AuthToken:  viper.GetString("gateways.twilio.auth_token"),
		FromNumber: viper.GetString("gateways.twilio.from_number"),
	}
}

func main() {
	g := viper.GetString("default_gateway")
	options := make([]string, len(gateways))
	i := 0
	for k := range gateways {
		options[i] = k
		i++
	}
	optionsString := strings.Join(options, ",")
	flag.StringVar(&g, "gateway", g, fmt.Sprintf("the gateway to send the SMS (options: %s)", optionsString))
	flag.StringVar(&g, "g", g, fmt.Sprintf("the gateway to send the SMS (shorthand) (options: %s)", optionsString))
	flag.Parse()

	var err error
	if g == "" {
		if viper.GetString("default_gateway") != "" {
			exitWithErrorMessage("gateway is required")
		}
		fmt.Println("no gateway was specified and no default is set")
		g, err = setDefaultGateway()
		if err != nil {
			exitWithError(err)
		}
		if g == "" {
			exitWithErrorMessage("gateway is required")
		}
	}

	gateway := gateways[g]
	if gateway == nil {
		exitWithErrorMessage("the gateway was not specified or could not be found: %s", g)
	}

	var recipient, text string
	argCount := len(os.Args)
	if argCount == 3 {
		recipient = os.Args[1]
		text = os.Args[2]
	} else if argCount == 4 || argCount == 5 {
		recipient = os.Args[argCount-2]
		text = os.Args[argCount-1]
	} else {
		flag.Usage()
		fmt.Println("qsms [recipient] [text]")
		exitWithErrorMessage("invalid argument count")
	}

	if gateway.RequiresConfiguration() {
		fmt.Printf("the %s gateway is not configured\n", g)
		if err = gateway.PromptForConfiguration(); err != nil {
			exitWithError(err)
		}
	}

	if err = gateway.Send(recipient, text); err != nil {
		exitWithError(err)
	}
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		exitWithError(err)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".qsms")
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err == nil {
		return
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		exitWithError(err)
	}

	// config file doesn't exist - create it
	f, err := os.Create(path.Join(home, ".qsms.json"))
	if err != nil {
		exitWithError(err)
	}
	if err = f.Close(); err != nil {
		exitWithError(err)
	}

	viper.Set("default_gateway", "")
	viper.Set("gateways.email.smtp_server", "")
	viper.Set("gateways.email.smtp_server_port", 0)
	viper.Set("gateways.email.sender", "")
	viper.Set("gateways.email.sender_password", "")
	viper.Set("gateways.email.mapping", map[string]string{"verizon": "vtext.com"})
	viper.Set("gateways.twilio.account_sid", "")
	viper.Set("gateways.twilio.auth_token", "")
	viper.Set("gateways.twilio.from_number", "")
	if err = viper.WriteConfig(); err != nil {
		exitWithError(err)
	}
}

func setDefaultGateway() (string, error) {
	fmt.Printf("would you like to set a default gateway now? (y/n): ")
	r := bufio.NewReader(os.Stdin)
	l, _, err := r.ReadLine()
	if err != nil {
		return "", err
	}

	ans := strings.ToLower(strings.TrimSpace(string(l)))
	if ans != "y" && ans != "yes" {
		return "", nil
	}

	options := make([]string, len(gateways))
	i := 0
	for k := range gateways {
		options[i] = k
		i++
	}

	fmt.Printf("select a gateway (options: %s): ", strings.Join(options, ","))
	l, _, err = r.ReadLine()
	if err != nil {
		return "", err
	}

	gateway := strings.ToLower(strings.TrimSpace(string(l)))
	if gateway == "" {
		return "", fmt.Errorf("gateway is required")
	} else if gateways[gateway] == nil {
		return "", fmt.Errorf("invalid option")
	}

	viper.Set("default_gateway", gateway)
	if err = viper.WriteConfig(); err != nil {
		return "", err
	}

	return gateway, nil
}

func exitWithErrorMessage(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}

func exitWithError(err error) {
	exitWithErrorMessage(err.Error())
}
