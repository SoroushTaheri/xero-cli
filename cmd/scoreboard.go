/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// scoreboardCmd represents the scoreboard command
var scoreboardCmd = &cobra.Command{
	Use:     "scoreboard",
	Aliases: []string{"sb"},
	Short:   "Show the scoreboard of the competition",
	Long:    `Show the scoreboard of the competition.`,
	Args:    cobra.ExactArgs(1),
	Run:     showScoreboard,
}

func init() {
	rootCmd.AddCommand(scoreboardCmd)
	scoreboardCmd.Flags().IntP("page", "p", 1, "Page of scoreboard")
}

func showScoreboard(cmd *cobra.Command, args []string) {
	problemPath := args[0]
	page, err := cmd.Flags().GetInt("page")
	if err != nil {
		fmt.Printf("failed to get pages from input command: %v\n", err)
		return
	}
	limit := 25
	start := (page - 1) * limit
	end := start + limit

	response, err := client.GetLeaderboard(problemPath, start, end)
	if err != nil {
		fmt.Printf("failed to get scoreboard data: %v\n", err)
		return
	}

	tableRows := pterm.TableData{
		{"Pos", "Team Name", "Captured", "Total Submissions", "Last Submission"},
	}

	tableRows = append(tableRows, []string{})
	for index, row := range response.Results {
		teamName := row.TeamName
		if !row.Individual {
			teamName = "[TEAM] " + teamName
		}

		tableRows = append(tableRows, []string{fmt.Sprint((page-1)*limit + (index + 1)), teamName, fmt.Sprint(row.Score), fmt.Sprint(len(row.SubmissionDates)), row.LastSubmission})
	}

	pterm.DefaultHeader.Println("Scoreboard")
	fmt.Printf("Total Records: %d\n\n", response.Total)
	pterm.DefaultTable.WithHasHeader().WithData(tableRows).Render()
}
