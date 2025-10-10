/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mdxabu/dodo/internal/scanner"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		scanner.Scan()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
