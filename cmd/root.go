/*
Copyright © 2025 Pay.jp
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/payjp/payjp-cli/internal/profiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "payjp-cli",
	Short:        "A CLI to help you integrate Pay.jp with your application",
	Long:         "The official command-line tool to interact with Pay.jp.",
	SilenceUsage: true,
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("profile-file-path", "", "file path (default is $HOME/.payjp-cli)")
	rootCmd.PersistentFlags().StringP("profile", "p", "default", "profile name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	profileFilePath := rootCmd.Flag("profile-file-path").Value.String()

	if profileFilePath != "" {
		profiles, err := profiles.LoadFromFile(profileFilePath)
		if err != nil {
			fmt.Println("Error loading profiles: ", err)
			os.Exit(1)
		}

		viper.Set("profiles", profiles)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".payjp-cli")

		profiles, err := profiles.LoadFromFile(home + "/.payjp-cli")
		if err != nil {
			fmt.Println("Error loading profiles: ", err)
			os.Exit(1)
		}

		viper.Set("profiles", profiles)
	}

	viper.SetDefault("BASE_URL", "https://pay.jp")
	viper.SetDefault("BASE_API_URL", "https://api.pay.jp")
	viper.SetDefault("GRPC_SERVER_ADDRESS", "https://grpc.pay.jp/cli")

	viper.AutomaticEnv() // read in environment variables that match

	profiles := viper.Get("profiles").(*profiles.Profiles)
	profileName := rootCmd.Flag("profile").Value.String()

	profile := profiles.LoadProfile(profileName)
	if profile != nil {
		if profile.BaseURL != "" {
			viper.Set("BASE_URL", profile.BaseURL)
		}
		if profile.BaseApiURL != "" {
			viper.Set("BASE_API_URL", profile.BaseApiURL)
		}
		if profile.GrpcServerAddress != "" {
			viper.Set("GRPC_SERVER_ADDRESS", profile.GrpcServerAddress)
		}
	}
}
