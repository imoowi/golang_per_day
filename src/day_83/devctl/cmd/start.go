/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var port int
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("server starting at :%d\n", port)
	},
}

func init() {
	serverCmd.AddCommand(startCmd)
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "listen port")
}
