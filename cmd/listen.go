/*
Copyright © 2025 Pay.jp
*/
package cmd

import (
	"fmt"

	pb "github.com/payjp/payjp-cli/gen/proto"
	"github.com/payjp/payjp-cli/internal/listen"
	"github.com/payjp/payjp-cli/internal/profiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen to events",
	Long:  `Listen to events from the PAY.JP. You can specify the events you want to listen to using the --events flag. By default, it listens to all events.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("listen called")

		profileName := cmd.Flag("profile").Value.String()
		allProfiles := viper.Get("profiles").(*profiles.Profiles)
		profile := allProfiles.GetProfile(profileName)
		if profile == nil {
			return fmt.Errorf("profile %s not found. you can create a profile using login command first.", profileName)
		}

		ctx := cmd.Context()
		address := viper.GetString("GRPC_SERVER_ADDRESS")

		err := listen.StartStream(ctx, address, &pb.ListenRequest{
			ApiKey: profile.TestModeSecretKey,
			Events: cmd.Flag("events").Value.String(),
		}, func(res *pb.ListenResponse) error {
			fmt.Printf("headers: %s \n", res.Headers)
			fmt.Printf("event received id: %s type: %s \n", res.PayjpEvent.Id, res.PayjpEvent.Type)
			// TODO: forwards the event to the local server
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")
	listenCmd.Flags().StringP("events", "e", "*", "events to listen to")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
