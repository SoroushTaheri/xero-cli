/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var challengeCmd = &cobra.Command{
	Use:     "challenge",
	Aliases: []string{"ch"},
	Short:   "Get the full list of challenges or view details of a challenge",
	Long:    `Get the full list of challenges or view details of a challenge`,
}

func init() {
	rootCmd.AddCommand(challengeCmd)
}
