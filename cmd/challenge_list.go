/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "Show a list of all challenges",
	Long:    `Show a list of all challenges`,
	Run:     list,
}

func init() {
	challengeCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) {
	if !client.IsLoggedIn() {
		fmt.Printf("❌ You are not logged in.\n\nTry logging in using: %q\n", "xero auth login")
		return
	}

	response, err := client.GetCompetition(config.Competition.Path)
	if err != nil {
		fmt.Printf("failed to get competition data: %v\n", err)
		return
	}

	// fmt.Println("❔ Challenges")
	pterm.DefaultHeader.Println("Challenges")

	problems := response.PhaseSet[0].Problems
	problemTitles := make([]pterm.BulletListItem, len(problems))

	for i, p := range problems {
		problemTitles[i] = pterm.BulletListItem{
			Level: 2,
			Text:  fmt.Sprintf("%s\t[%s]", p.Problem.Title, p.Problem.Path),
		}
	}

	pterm.DefaultBulletList.WithItems(problemTitles).Render()
}
