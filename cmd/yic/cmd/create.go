package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/youritcity/go-sdk/youritcity"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <role>",
	Short: "Create an app token to yourITcity service",
	Long: `Create an app token to yourITcity service. 
	
It is like a login but without the need to validate the token by email.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appToken, ok := os.LookupEnv("YOURITCITY_APPTOKEN")
		failOnCond(!ok, fmt.Sprintf("Missing token, please set %s with 'yic login' or 'yic signup' command.", EnvYourITcityToken))

		conn := youritcity.NewConnection()
		_, err := conn.Renew(youritcity.App{Token: appToken})
		failOnError(err, "Can't renew session token. Please validate this token or the token has been revoked. You can a new one by using the login subcommand.")

		auth := struct {
			Name string `json:"name,omitempty"`
			Type string `json:"type"`
		}{
			Type: args[0],
			Name: "yic cli",
		}
		body, err := json.Marshal(auth)
		failOnError(err, "can't craft request")
		resp, err := conn.Client().Post(youritcity.YourITcityURI+"/account/token", "application/json", bytes.NewBuffer(body))
		failOnError(err, "can't create token")
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			failOnCond(true, fmt.Sprintf("%s: %s", resp.Status, b))
		}
		fmt.Printf("%s\n", b)
	},
}

func init() {
	tokenCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
