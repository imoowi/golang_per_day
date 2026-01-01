/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"codee_jun/internal/config"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var env string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "codee_jun",
	Short: "codee_jun",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		//初始化配置
		return config.Load(cfgFile, env)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.day56.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is configs/confg.yaml)")
	rootCmd.PersistentFlags().StringVar(&env, "env", "dev", "runtime environment: dev/prod")
}
