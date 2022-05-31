/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"xero-cli/pkg/paricheh"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "",
	Long:  ``,
	Run:   submit,
	Args:  cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(submitCmd)
}

type ParichehRequest struct {
	Flag string `json:"flag"`
}

func submit(cmd *cobra.Command, args []string) {
	problemPath := args[0]
	flag := args[1]

	validFlag, err := paricheh.SendSubmittedFlag(problemPath, client.AccessToken, flag)
	if err != nil {
		fmt.Printf("❌ Failed to submit flag: %v\n", err)
		return
	}

	if !validFlag {
		fmt.Println("❌ Invalid flag")
		return
	}

	fmt.Println("✅ Correct! You did it!")

}
