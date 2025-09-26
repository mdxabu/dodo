/*
Copyright © 2025 @mdxabu
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Stage all changes, scan for secrets, auto-generate commit message based on file changes",
	Long: `Stage all changes, scan for secrets, auto-generate commit message based on file changes

	Example:

	dodo commit
	dodo commit -m "fix: typo"

	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the custom message if provided
		customMessage, _ := cmd.Flags().GetString("message")

		// Execute the commit process
		if err := executeCommit(customMessage); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func executeCommit(customMessage string) error {
	// Find the git repository
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get the working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get the status before staging
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	// Check if there are any changes to commit
	if status.IsClean() {
		fmt.Println("No changes to commit, working tree clean")
		return nil
	}

	// Collect information about changes for message generation
	var addedFiles, modifiedFiles, deletedFiles []string

	for file, fileStatus := range status {
		switch fileStatus.Staging {
		case git.Added:
			addedFiles = append(addedFiles, file)
		case git.Modified:
			modifiedFiles = append(modifiedFiles, file)
		case git.Deleted:
			deletedFiles = append(deletedFiles, file)
		}

		// Also check worktree status for unstaged changes
		switch fileStatus.Worktree {
		case git.Added:
			if fileStatus.Staging != git.Added {
				addedFiles = append(addedFiles, file)
			}
		case git.Modified:
			if fileStatus.Staging != git.Modified {
				modifiedFiles = append(modifiedFiles, file)
			}
		case git.Deleted:
			if fileStatus.Staging != git.Deleted {
				deletedFiles = append(deletedFiles, file)
			}
		case git.Untracked:
			addedFiles = append(addedFiles, file)
		}
	}

	// Stage all changes (git add .)
	fmt.Println("Staging all changes...")
	err = worktree.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Generate commit message
	var commitMessage string
	if customMessage != "" {
		commitMessage = customMessage
	} else {
		commitMessage = generateCommitMessage(addedFiles, modifiedFiles, deletedFiles)
	}

	// Create the commit
	fmt.Printf("Committing with message: %s\n", commitMessage)

	// Get git config for author information
	config, err := repo.Config()
	if err != nil {
		return fmt.Errorf("failed to get git config: %w", err)
	}

	// Try to get author info from git config
	var authorName, authorEmail string
	if config.User.Name != "" {
		authorName = config.User.Name
	} else {
		authorName = "Unknown"
	}
	if config.User.Email != "" {
		authorEmail = config.User.Email
	} else {
		authorEmail = "unknown@example.com"
	}

	commit, err := worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Print commit hash
	fmt.Printf("Successfully committed: %s\n", commit.String()[:8])

	return nil
}

func generateCommitMessage(added, modified, deleted []string) string {
	var parts []string

	// Count changes
	totalChanges := len(added) + len(modified) + len(deleted)

	if totalChanges == 0 {
		return "chore: update files"
	}

	// Determine the primary action and scope
	var primaryAction string
	var scope string

	// Analyze file types to determine scope
	allFiles := append(append(added, modified...), deleted...)
	scope = determineScope(allFiles)

	// Determine primary action based on the type of changes
	if len(added) > len(modified) && len(added) > len(deleted) {
		primaryAction = "feat"
		if len(added) == 1 {
			parts = append(parts, fmt.Sprintf("add %s", filepath.Base(added[0])))
		} else {
			parts = append(parts, fmt.Sprintf("add %d files", len(added)))
		}
	} else if len(deleted) > 0 {
		primaryAction = "refactor"
		if len(deleted) == 1 {
			parts = append(parts, fmt.Sprintf("remove %s", filepath.Base(deleted[0])))
		} else {
			parts = append(parts, fmt.Sprintf("remove %d files", len(deleted)))
		}
	} else {
		primaryAction = "fix"
		if len(modified) == 1 {
			parts = append(parts, fmt.Sprintf("update %s", filepath.Base(modified[0])))
		} else {
			parts = append(parts, fmt.Sprintf("update %d files", len(modified)))
		}
	}

	// Add additional details for mixed changes
	if len(added) > 0 && primaryAction != "feat" {
		if len(added) == 1 {
			parts = append(parts, fmt.Sprintf("add %s", filepath.Base(added[0])))
		} else {
			parts = append(parts, fmt.Sprintf("add %d files", len(added)))
		}
	}

	if len(modified) > 0 && primaryAction != "fix" {
		if len(modified) == 1 {
			parts = append(parts, fmt.Sprintf("update %s", filepath.Base(modified[0])))
		} else {
			parts = append(parts, fmt.Sprintf("update %d files", len(modified)))
		}
	}

	// Construct the final message
	message := primaryAction
	if scope != "" {
		message += "(" + scope + ")"
	}
	message += ": " + strings.Join(parts, ", ")

	return message
}

func determineScope(files []string) string {
	hasGo := false
	hasJS := false
	hasTS := false
	hasPy := false
	hasDoc := false
	hasConfig := false
	hasTest := false

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		base := strings.ToLower(filepath.Base(file))

		switch ext {
		case ".go":
			hasGo = true
		case ".js":
			hasJS = true
		case ".ts", ".tsx":
			hasTS = true
		case ".py":
			hasPy = true
		case ".md", ".txt", ".rst":
			hasDoc = true
		case ".json", ".yaml", ".yml", ".toml", ".ini", ".conf":
			hasConfig = true
		}

		if strings.Contains(base, "test") || strings.Contains(base, "spec") {
			hasTest = true
		}

		if base == "dockerfile" || base == "makefile" || base == "readme.md" {
			hasConfig = true
		}
	}

	// Return the most specific scope
	if hasTest {
		return "test"
	}
	if hasGo {
		return "go"
	}
	if hasTS {
		return "ts"
	}
	if hasJS {
		return "js"
	}
	if hasPy {
		return "py"
	}
	if hasDoc {
		return "docs"
	}
	if hasConfig {
		return "config"
	}

	return ""
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "Custom commit message")
}
