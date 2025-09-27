/*
Copyright © 2025 @mdxabu
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes/staged commits to the remote repository",
	Long: `Push changes/staged commits to the remote repository.

	Example:

	dodo push
	dodo push --force
	dodo push --set-upstream origin main

	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		force, _ := cmd.Flags().GetBool("force")
		setUpstream, _ := cmd.Flags().GetString("set-upstream")

		// Execute the push process
		if err := executePush(force, setUpstream, args); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func executePush(force bool, setUpstream string, args []string) error {
	// Find the git repository
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get the current HEAD reference
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Get current branch name
	branchName := head.Name().Short()
	fmt.Printf("Current branch: %s\n", branchName)

	// Determine remote and branch to push to
	var remoteName, refSpec string
	if len(args) >= 2 {
		// User specified remote and branch: dodo push origin main
		remoteName = args[0]
		refSpec = fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, args[1])
	} else if len(args) == 1 {
		// User specified only remote: dodo push origin
		remoteName = args[0]
		refSpec = fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)
	} else {
		// No args, use default remote (usually "origin")
		remoteName = "origin"
		refSpec = fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)
	}

	// Handle set-upstream option
	if setUpstream != "" {
		remoteName = setUpstream
		if len(args) > 0 {
			refSpec = fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, args[0])
		} else {
			refSpec = fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)
		}
	}

	fmt.Printf("Pushing to remote '%s'...\n", remoteName)

	// Verify the remote exists
	_, err = repo.Remote(remoteName)
	if err != nil {
		return fmt.Errorf("failed to get remote '%s': %w", remoteName, err)
	}

	// Create push options
	pushOptions := &git.PushOptions{
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec(refSpec),
		},
		Force: force,
	}

	// Perform the push
	err = repo.Push(pushOptions)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("Everything up-to-date")
			return nil
		}
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Printf("Successfully pushed %s to %s\n", branchName, remoteName)

	// If set-upstream was used, update the local config
	if setUpstream != "" && len(args) > 0 {
		fmt.Printf("Branch '%s' set up to track remote branch '%s' from '%s'\n", branchName, args[0], setUpstream)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolP("force", "f", false, "Force push (use with caution)")
	pushCmd.Flags().StringP("set-upstream", "u", "", "Set upstream remote and branch")
}
