/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rulesCmd represents the rules command
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "View CTF rules",
	Long:  `View CTF rules`,
	Run:   rules,
}

func init() {
	rootCmd.AddCommand(rulesCmd)
}

func rules(cmd *cobra.Command, args []string) {
	response, err := client.GetCompetition(config.Competition.Path)
	if err != nil {
		fmt.Printf("failed to get competition data: %v\n", err)
		return
	}

	fmt.Println("🚧 XeroCTF Rules")
	fmt.Println(response.Rules)
}
