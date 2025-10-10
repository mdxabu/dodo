/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("scan called")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
