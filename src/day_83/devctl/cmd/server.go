/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the development server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called pre")
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("server called pre run error")
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called persistent pre")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
