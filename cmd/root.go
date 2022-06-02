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

var rootCmd = &cobra.Command{
	Use:   "xero",
	Short: "Command-line interface for XeroCTF 2022",
	Long:  `Command-line interface for XeroCTF 2022.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
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
