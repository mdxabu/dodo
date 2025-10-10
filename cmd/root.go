/*
Copyright © 2025 @mdxabu
*/
package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dodo",
	Short: "",
	Long:  `dodo is a Go-based command-line tool designed to make Git safer and easier for developers by solving pains...`,
	Run: func(cmd *cobra.Command, args []string) {

		asciibanner := `
      $$\                 $$\           
      $$ |                $$ |          
 $$$$$$$ | $$$$$$\   $$$$$$$ | $$$$$$\  
$$  __$$ |$$  __$$\ $$  __$$ |$$  __$$\ 
$$ /  $$ |$$ /  $$ |$$ /  $$ |$$ /  $$ |
$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |
\$$$$$$$ |\$$$$$$  |\$$$$$$$ |\$$$$$$  |
 \_______| \______/  \_______| \______/ 
                                        
                                        
                                        
		`

		fmt.Println(asciibanner)
        fmt.Println("Welcome to dodo! Use 'dodo --help' to see available commands.")
    },

}


func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
