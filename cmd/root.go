/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	roboepicsClient "xero-cli/pkg/client"
	xeroConfig "xero-cli/pkg/config"
)

var (
	client *roboepicsClient.Client
	config *xeroConfig.Config
	err    error
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xero",
	Short: "Command-line interface for XeroCTF 2022",
	Long:  `Command-line interface for XeroCTF 2022.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.xero-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	config, err = xeroConfig.New()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize config: %v", err))
	}

	if len(config.AccessToken)*len(config.RefreshToken) == 0 {
		client = roboepicsClient.New()
	} else {
		client, err = roboepicsClient.NewWithToken(config.AccessToken, config.RefreshToken, false)
		if err != nil {
			panic(fmt.Sprintf("failed to initialize RoboEpics client: %v", err))
		}

		config.SetTokens(client.AccessToken, client.RefreshToken)
	}
}
