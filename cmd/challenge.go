/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// challengeCmd represents the challenge command
var challengeCmd = &cobra.Command{
	Use:     "challenge",
	Aliases: []string{"ch"},
	Short:   "Get the full list of challenges or view details of a challenge",
	Long:    `Get the full list of challenges or view details of a challenge`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("challenge called")
	},
}

func init() {
	rootCmd.AddCommand(challengeCmd)
}
