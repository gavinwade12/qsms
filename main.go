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
	senderPass, _ := base32.StdEncoding.DecodeString(viper.GetString("gateways.email.sender_password"))
	gateways["email"] = EmailGateway{
		SMTPServer:     viper.GetString("gateways.email.smtp_server"),
		Sender:         viper.GetString("gateways.email.sender"),
		SenderPassword: string(senderPass),
	}
}

func main() {
	initConfig()

	g := viper.GetString("default_gateway")
	flag.StringVar(&g, "gateway", g, "the gateway to send the SMS")
	flag.StringVar(&g, "g", g, "the gateway to send the SMS (shorthand)")
	flag.Parse()

	if len(os.Args) < 3 {
		flag.Usage()
		fmt.Println("qsms [recipient] [number]")
		return
	}

	if g == "" {
		if viper.GetString("default_gateway") != "" {
			exitWithErrorMessage("gateway is required")
		}
		fmt.Println("no gateway was specified and no default is set")
		fmt.Printf("would you like to set a default gateway now? (y/n): ")
		r := bufio.NewReader(os.Stdin)
		l, _, err := r.ReadLine()
		if err != nil {
			exitWithError(err)
		}

		ans := strings.ToLower(strings.TrimSpace(string(l)))
		if ans != "y" && ans != "yes" {
			exitWithErrorMessage("a gateway is required")
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
			exitWithError(err)
		}

		ans = strings.ToLower(strings.TrimSpace(string(l)))
		if ans == "" {
			exitWithErrorMessage("gateway is required")
		} else if gateways[ans] == nil {
			exitWithErrorMessage("invalid option")
		}

		viper.Set("default_gateway", ans)
		if err = viper.WriteConfig(); err != nil {
			exitWithError(err)
		}

		g = ans
	}

	gateway := gateways[g]
	if gateway == nil {
		exitWithErrorMessage("the gateway was not specified or could not be found: %s", g)
	}

	if err := gateway.Send(os.Args[1], os.Args[2]); err != nil {
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

	viper.AutomaticEnv() // read in environment variables that match

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
	if err = viper.WriteConfig(); err != nil {
		exitWithError(err)
	}
}

func exitWithErrorMessage(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}

func exitWithError(err error) {
	exitWithErrorMessage(err.Error())
}
