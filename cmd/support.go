/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var supportCmd = &cobra.Command{
	Use:   "support",
	Short: "Link to our discord server",
	Long:  `Link to our discord server`,
	Run:   support,
}

func init() {
	rootCmd.AddCommand(supportCmd)
}

func support(cmd *cobra.Command, args []string) {
	fmt.Println("Join our Discord server:")
	fmt.Println("https://discord.gg/8DMfjmn6gc")
}
