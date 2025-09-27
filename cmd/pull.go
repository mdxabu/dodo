/*
Copyright © 2025 @mdxabu
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from the remote repository",
	Long: `Pull changes from the remote repository.

	Example:

	dodo pull
	dodo pull origin
	dodo pull origin main
	dodo pull --rebase

	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		rebase, _ := cmd.Flags().GetBool("rebase")
		force, _ := cmd.Flags().GetBool("force")

		// Execute the pull process
		if err := executePull(rebase, force, args); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func executePull(rebase bool, force bool, args []string) error {
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

	// Determine remote and branch to pull from
	var remoteName string
	if len(args) >= 1 {
		// User specified remote: dodo pull origin [branch]
		remoteName = args[0]
	} else {
		// No args, use default remote (usually "origin")
		remoteName = "origin"
	}

	fmt.Printf("Pulling from remote '%s'...\n", remoteName)

	// Verify the remote exists
	_, err = repo.Remote(remoteName)
	if err != nil {
		return fmt.Errorf("failed to get remote '%s': %w", remoteName, err)
	}

	// Get the working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Create pull options
	pullOptions := &git.PullOptions{
		RemoteName: remoteName,
		Force:      force,
	}

	// If a specific branch was specified, set the reference name
	if len(args) >= 2 {
		pullOptions.ReferenceName = head.Name()
	}

	// Note: go-git doesn't have built-in rebase support for pull
	// For now, we'll mention this limitation but still pull normally
	if rebase {
		fmt.Println("Note: Rebase option noted, but go-git performs merge by default")
	}

	// Perform the pull
	err = worktree.Pull(pullOptions)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("Already up to date")
			return nil
		}
		return fmt.Errorf("failed to pull: %w", err)
	}

	fmt.Printf("Successfully pulled from %s to %s\n", remoteName, branchName)

	return nil
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().BoolP("rebase", "r", false, "Rebase the current branch on top of the upstream branch after fetching")
	pullCmd.Flags().BoolP("force", "f", false, "Force pull (use with caution)")
}
