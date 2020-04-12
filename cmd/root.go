package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var gateway map[string]string

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "qsms",
	Short: "qsms is used for quickly sending text messages",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qsms.json)")
}

func initConfig() {
	var home string
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		var err error
		home, err = homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".qsms")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			log.Fatal(err)
		}

		f, err := os.Create(path.Join(home, ".qsms.json"))
		if err != nil {
			log.Fatal(err)
		}
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}

		viper.Set("gateway", map[string]string{"verizon": "vtext.com"})
		if err = viper.WriteConfig(); err != nil {
			log.Fatal(err)
		}
	}

	gateway = viper.GetStringMapString("gateway")
}
