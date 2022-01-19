/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/criscola/cloudflare-ddns/ddns"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"

	"github.com/spf13/cobra"
)

var (
	logger   *zap.Logger
	cfClient *cloudflare.API
	config   ddns.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudflare-ddns",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("config")         // name of config file (without extension)
		viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
		viper.AddConfigPath(".")              // optionally look for config in the working directory
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				fmt.Println("Config file not found.")
			} else {
				fmt.Println("Error during config read.")
			}
		}
		err := viper.Unmarshal(&config)
		if err != nil {
			// Unmarshalling error
			fmt.Println("Error during config unmarshalling.")
		}
		fmt.Println("Yay")
		fmt.Println(config)
	},
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudflare-ddns.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
