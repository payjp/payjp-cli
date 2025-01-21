/*
Copyright © 2025 Pay.jp
*/
package cmd

import (
	"fmt"

	"github.com/payjp/payjp-cli/internal/listen"
	"github.com/payjp/payjp-cli/internal/payjp"
	"github.com/payjp/payjp-cli/internal/profiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("listen called")

		profileName := cmd.Flag("profile").Value.String()
		allProfiles := viper.Get("profiles").(*profiles.Profiles)
		profile := allProfiles.GetProfile(profileName)
		if profile == nil {
			return fmt.Errorf("profile %s not found. you can create a profile using login command first.", profileName)
		}

		ctx := cmd.Context()
		apiClient, err := payjp.NewClient(viper.GetString("BASE_API_URL"), profile.TestModeSecretKey)
		if err != nil {
			return err
		}

		session, err := listen.CreateCliSession(ctx, apiClient)
		if err != nil {
			return err
		}

		fmt.Println("Session: ", session)

		// TODO: grpc 接続

		fmt.Println("Listening finished")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
