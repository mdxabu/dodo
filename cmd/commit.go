/*
Copyright © 2025 @mdxabu

*/
package cmd

import (
	"github.com/go-git/go-git/v5"
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
		
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().BoolP("message", "m", false, "Enter the commit message")
}
