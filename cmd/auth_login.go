/*
Copyright © 2022 Soroush Taheri soroushtgh@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login using your RoboEpics credentials",
	Long:  `Login using your RoboEpics credentials to be able to execute protected/user-based commands. (submissions, challenge descriptions etc.)`,
	Run:   login,
	// Args:  cobra.ExactArgs(2),
}

func init() {
	authCmd.AddCommand(loginCmd)
}

func login(cmd *cobra.Command, args []string) {
	defer config.Sync()

	username, err := promptUsername()
	if err != nil {
		// fmt.Printf("❌ Username prompt failed: %v\n", err)
		return
	}

	password, err := promptPassword()
	if err != nil {
		// fmt.Printf("❌ Password prompt failed: %v\n", err)
		return
	}

	if err := client.Login(username, password); err != nil {
		fmt.Printf("❌ Login Error: %v\n", err)
		return
	}

	config.SetTokens(client.Creds())
	fmt.Printf("✅ Successfully logged in as: %s\n", username)
}

func promptPassword() (string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . | bold }}: ",
		Valid:   "{{ . | bold }}: ",
		Invalid: "{{ . | bold }}: ",
		Success: "{{ . | bold }}: ",
	}

	prompt := promptui.Prompt{
		Label: "Password",
		// HideEntered: true,
		Templates: templates,
		Mask:      '*',
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func promptUsername() (string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . | bold }}: ",
		Valid:   "{{ . | bold }}: ",
		Invalid: "{{ . | bold }}: ",
		Success: "{{ . | bold }}: ",
	}

	prompt := promptui.Prompt{
		Label:     "Username/Email",
		Templates: templates,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}
