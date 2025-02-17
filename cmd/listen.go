/*
Copyright © 2025 Pay.jp
*/
package cmd

import (
	"fmt"
	"log"
	"net/url"
	"path"

	pb "github.com/payjp/payjp-cli/gen/proto"
	"github.com/payjp/payjp-cli/internal/ansi"
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
		profileName := cmd.Flag("profile").Value.String()
		allProfiles := viper.Get("profiles").(*profiles.Profiles)
		profile := allProfiles.LoadProfile(profileName)
		if profile == nil {
			return fmt.Errorf("profile %s not found. you can create a profile using login command first.", profileName)
		}

		ctx := cmd.Context()

		forwardTo := cmd.Flag("forward-to").Value.String()
		skipVerify, err := cmd.Flags().GetBool("skip-verify")
		if err != nil {
			return err
		}

		eventForwarder, err := listen.NewEventForwarder(forwardTo, skipVerify)
		if err != nil {
			return err
		}
		address := viper.GetString("GRPC_SERVER_ADDRESS")

		listener := listen.NewListener(address)
		err = listener.StartListen(ctx, &pb.ListenRequest{
			ApiKey: profile.TestModeSecretKey,
			Events: cmd.Flag("events").Value.String(),
		}, func(res *pb.PayjpEventResponse) error {
			eventURL, err := url.Parse(viper.GetString("BASE_URL"))
			if err != nil {
				return err
			}
			eventURL.Path = path.Join(eventURL.Path, "d", "events", res.PayjpEvent.Id)
			log.Printf("--> %s [%s]\n", res.PayjpEvent.Type, ansi.Link(eventURL.String(), res.PayjpEvent.Id))

			err = eventForwarder.ForwardEvent(res)
			if err != nil {
				return err
			}

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

	listenCmd.Flags().StringP("events", "e", "*", "comma separated list of events to listen to, e.g. charge.created,charge.updated")
	listenCmd.Flags().StringP("forward-to", "f", "", "forward events to this localhost port number or url, e.g. 3000, https://example.com/hook")
	listenCmd.Flags().Bool("skip-verify", false, "skip the verification of the SSL certificate of the forward-to url")
}
