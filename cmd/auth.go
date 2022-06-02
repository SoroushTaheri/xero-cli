/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Auth-related commands (login, status, ...)",
	Long:  `Authentication-related commands let you login to your RoboEpics account in order to be able to use user-specific commands (eg. submitting a flag)`,
}

func init() {
	rootCmd.AddCommand(authCmd)
}
