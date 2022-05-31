/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Check your login information",
	Long:    `Check if you're logged in or not and who you're logged in as.`,
	Run:     status,
}

func init() {
	authCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) {
	if len(client.AccessToken) == 0 {
		fmt.Printf("❌ You are not logged in.\n\nTry logging in using: %q\n", "xero auth login")
		return
	}
	client.UpdateProfileFromToken()
	fmt.Printf("Logged in as \"%s/%s\"\n", client.Profile.Email, client.Profile.Username)
}
