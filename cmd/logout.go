/*
Copyright © 2025 Pay.jp
*/
package cmd

import (
	"fmt"

	"github.com/payjp/payjp-cli/internal/profiles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the PAY.JP CLI",
	Long:  "Use this command to log out from the PAY.JP CLI.  You can log out from a specific profile or all profiles.",
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := cmd.Flag("profile").Value.String()
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}

		allProfiles := viper.Get("profiles").(*profiles.Profiles)
		if all {
			allProfiles.Profiles = make(map[string]*profiles.Profile)
			err = allProfiles.Persist()
			if err != nil {
				fmt.Println("Error saving profile: ", err)
				return err
			}

			fmt.Println("All profiles logged out.")
			return nil
		}

		allProfiles.RemoveProfile(profileName)
		err = allProfiles.Persist()
		if err != nil {
			fmt.Println("Error saving profile: ", err)
			return err
		}

		fmt.Printf("Profile %s logged out.\n", profileName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
	logoutCmd.Flags().BoolP("all", "a", false, "Clear credentials for all projects you are currently logged into.")
}
