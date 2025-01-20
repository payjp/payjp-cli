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

var profileFilePath string

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

	rootCmd.PersistentFlags().StringVarP(&profileFilePath, "profile-file-path", "f", "", "profile file path (default is $HOME/.payjp-cli)")
	rootCmd.PersistentFlags().StringP("profile", "p", "default", "profile name")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

	viper.AutomaticEnv() // read in environment variables that match
}
