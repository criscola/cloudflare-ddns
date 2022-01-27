/*
Copyright © 2022 Cristiano Colangelo (criscola)

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/criscola/cloudflare-ddns/ddns"
	"github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
)

var (
	logger   *zap.SugaredLogger
	cfClient *cloudflare.API
	config   ddns.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudflare-ddns",
	Short: "A Dynamic DNS for Cloudflare that allows your zones to be constantly up-to-date with your constantly changing public IP",
	Long:  `A Dynamic DNS for Cloudflare that allows your zones to be constantly up-to-date with your constantly changing public IP`,
	Run: func(cmd *cobra.Command, args []string) {
		client := ddns.Client{API: cfClient, Logger: logger, Config: config}
		client.Start()
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
	initLogger()
	initConfig()
	initClient()
	initProbes()

	// TODO: CLI Flags e.g. debug mode
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudflare-ddns.yaml)")
}

func initLogger() {
	temp, _ := zap.NewDevelopment()
	logger = temp.Sugar()
}

func initConfig() {
	viper.SetConfigName("config")                 // name of config file (without extension)
	viper.SetConfigType("yaml")                   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/cloudflare-ddns")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.cloudflare-ddns") // call multiple times to add many search paths
	viper.AddConfigPath(".")                      // optionally look for config in the working directory
	// TODO: Error handling, config validation, dependency injection
	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			logger.Fatal("config file not found.")
		} else {
			logger.Fatalf("error during config read: %s", err)
		}
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		// Unmarshalling error
		logger.Fatalf("error during config unmarshalling: %s", err)
	}
	logger.Debugf("config loaded successfully from: %s", viper.ConfigFileUsed())
}

func initClient() {
	// Construct a new API object using a global API key
	var err error
	cfClient, err = cloudflare.NewWithAPIToken(config.ApiToken)
	if err != nil {
		logger.Fatal(err)
	}
	// Test if connection works and can list DNS zones
	_, err = cfClient.ListZones(context.Background())
	if err != nil {
		logger.Fatalf("cannot list DNS zones: %s", err)
	}
	logger.Debug("cfClient initialized successfully")
}

func initProbes() {
	health := healthcheck.NewHandler()
	// If app cannot reach 1.1.1.1, readiness probe will fail
	checkPublicIpResolver := func() error {
		resp, err := http.Get(ddns.PublicIpPage)
		if err != nil {
			return err
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("public IP resolver (%s) returned status code %d",
				ddns.PublicIpPage, resp.StatusCode)
		}
		return nil
	}
	health.AddReadinessCheck("public-ip-resolver-webpage", checkPublicIpResolver)
	go func() {
		err := http.ListenAndServe("0.0.0.0:8086", health)
		if err != nil {
			logger.Fatalf("cannot start webserver on port 8086")
		}
		logger.Debug("probes started successfully")
	}()
}
