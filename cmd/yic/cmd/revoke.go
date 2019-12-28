package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fredericalix/yic_go_sdk/youritcity"

	"github.com/spf13/cobra"
)

// revokeCmd represents the delete command
var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke an application token",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appToken, ok := os.LookupEnv("YOURITCITY_APPTOKEN")
		failOnCond(!ok, fmt.Sprintf("Missing token, please set %s with 'yic login' or 'yic signup' command.", EnvYourITcityToken))

		conn := youritcity.NewConnection()
		_, err := conn.Renew(youritcity.App{Token: appToken})
		failOnError(err, "Can't renew session token. Please validate this token or the token has been revoked. You can a new one by using the login subcommand.")

		client := conn.Client()
		req, err := http.NewRequest(http.MethodDelete, youritcity.YourITcityURI+"/account/token/"+args[0], nil)
		failOnError(err, "can't build revoke request")
		resp, err := client.Do(req)
		failOnError(err, "can't delete request")
		if resp.StatusCode != http.StatusOK {
			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)
			failOnCond(true, fmt.Sprintf("%s: %s", resp.Status, b))
		}
	},
}

func init() {
	tokenCmd.AddCommand(revokeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
