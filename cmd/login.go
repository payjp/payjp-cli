/*
Copyright © 2025 Pay.jp
*/
package cmd

import (
	"fmt"

	"github.com/payjp/payjp-cli/internal/login"
	"github.com/payjp/payjp-cli/internal/payjp"
	"github.com/payjp/payjp-cli/internal/profiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and configure the PAY.JP CLI profile",
	Long:  "Use this command to authenticate and configure your PAY.JP CLI profile.  This command will guide you through the authentication process and set up your profile with the necessary credentials to interact with the PAY.JP API.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		dashboardClient, err := payjp.NewClient(viper.GetString("BASE_URL"), "")
		if err != nil {
			return err
		}

		authResponse, err := login.CallAuth(ctx, dashboardClient)
		if err != nil {
			return err
		}

		fmt.Printf("Your pairing code is: %s\n", authResponse.PairingCode)
		fmt.Printf("To authenticate with PAY.JP, please access to: %s\n", authResponse.BrowserURL)
		fmt.Println("Waiting for confirmation...")

		pollingCh := make(chan *login.AuthResult)
		go login.PollingAuthResult(ctx, dashboardClient, authResponse.PollURL, pollingCh)
		authResult := <-pollingCh
		if authResult.Err != nil {
			return authResult.Err
		}

		profileName := cmd.Flag("profile").Value.String()
		fmt.Printf(
			"Successfully authenticated! The PAY.JP CLI %s profile is configured for %s with account id %s\n",
			profileName,
			authResult.AccountDisplayName,
			authResult.AccountID,
		)

		loggedInProfile := &profiles.Profile{
			Name:              profileName,
			TestModeSecretKey: authResult.TestModeSecretKey,
			BaseURL:           viper.GetString("BASE_URL"),
			GrpcServerAddress: viper.GetString("GRPC_SERVER_ADDRESS"),
		}

		allProfiles := viper.Get("profiles").(*profiles.Profiles)
		allProfiles.AddProfile(loggedInProfile)
		err = allProfiles.Persist()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
