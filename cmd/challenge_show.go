/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"
	"strings"

	roboepicsClient "xero-cli/pkg/client"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a challenge",
	Long:  `Show details of a challenge`,
	Run:   show,
	Args:  cobra.ExactArgs(1),
}

func init() {
	challengeCmd.AddCommand(showCmd)
}

func show(cmd *cobra.Command, args []string) {
	if !client.IsLoggedIn() {
		fmt.Printf("❌ You are not logged in.\n\nTry logging in using: %q\n", "xero auth login")
		return
	}

	problemPath := args[0]

	response, err := client.GetProblem(problemPath)
	if err != nil {
		fmt.Printf("failed to get problem data: %v\ncheck your input\n", err)
		return
	}

	textsResponse, err := client.GetProblemTexts(problemPath)
	if err != nil {
		fmt.Printf("failed to get problem texts: %v\n", err)
		return
	}

	texts := make([]roboepicsClient.ProblemText, len(textsResponse))

	for index, problemTextRef := range textsResponse {
		textResponse, err := client.GetProblemText(fmt.Sprint(problemTextRef.ID))
		if err != nil {
			fmt.Printf("failed to get problem text: %v\n", err)
			return
		}

		texts[index] = textResponse
	}

	pterm.DefaultHeader.Println(response.Title)
	fmt.Println()

	for _, text := range texts {
		pterm.DefaultSection.Println(text.Title)
		for _, p := range strings.Split(text.Text, "\n\n") {
			pterm.DefaultParagraph.WithMaxWidth(150).Println(p)
			fmt.Printf("\n")
		}
	}
}
