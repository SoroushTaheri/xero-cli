/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Auth-related commands (login, status, ...)",
	Long:  `Authentication-related commands let you login to your RoboEpics account in order to be able to use user-specific commands (eg. submitting a flag)`,
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("auth called")
	// 	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
