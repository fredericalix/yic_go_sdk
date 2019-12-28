package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fredericalix/yic_go_sdk/youritcity"
	uuid "github.com/satori/go.uuid"

	"github.com/spf13/cobra"
)

func failOnError(err error, why string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", why, err)
		os.Exit(1)
	}
}

func failOnCond(nok bool, why string) {
	if nok {
		fmt.Fprintln(os.Stderr, why)
		os.Exit(1)
	}
}

type AppToken struct {
	Token      string `json:"app_token,omitempty"`
	ValidToken string `json:"validation_token,omitempty"`
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Roles      Roles  `json:"roles,omitempty"`

	AID uuid.UUID `json:"account_id,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}

// Roles is a list of name with their authorization (r: read, w: write)
type Roles map[string]string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your authentification token",
	Run: func(cmd *cobra.Command, args []string) {
		appToken, ok := os.LookupEnv("YOURITCITY_APPTOKEN")
		failOnCond(!ok, fmt.Sprintf("Missing token, please set %s with 'yic login' or 'yic signup' command.", EnvYourITcityToken))

		conn := youritcity.NewConnection()
		_, err := conn.Renew(youritcity.App{Token: appToken})
		failOnError(err, "Can't renew session token. Please validate this token or the token has been revoked. You can a new one by using the login subcommand.")

		client := conn.Client()

		resp, err := client.Get(youritcity.YourITcityURI + "/account/token")
		failOnError(err, "can't get list of token")
		failOnCond(resp.StatusCode != http.StatusOK, "responce status: "+resp.Status)
		defer resp.Body.Close()

		var tokens []struct {
			Token     string `json:"app_token,omitempty"`
			ValidLink string `json:"validation_link,omitempty"`
			Name      string `json:"name,omitempty"`
			Type      string `json:"type,omitempty"`

			Roles map[string]string `json:"roles,omitempty"`

			CreatedAt time.Time `json:"created_at,omitempty"`
			UpdatedAt time.Time `json:"updated_at,omitempty"`
			ExpiredAt time.Time `json:"expired_at,omitempty"`
		}
		err = json.NewDecoder(resp.Body).Decode(&tokens)
		failOnError(err, "can't parse responce")

		var tokenLen, typeLen, nameLen int
		for _, t := range tokens {
			if len(t.Token) > tokenLen {
				tokenLen = len(t.Token)
			}
			if len(t.Type) > typeLen {
				typeLen = len(t.Type)
			}
			if len(t.Name) > nameLen {
				nameLen = len(t.Name)
			}
		}

		fmt.Printf("%- [1]*s  %s  %- [4]*s  %- [6]*s  %- 10s  %- 10s  %-10s  %s\n",
			tokenLen, "Tokens", "V", nameLen, "Name", typeLen, "Type", "Created_At", "Updated_At", "Expired_At", "Validation_Link")
		for _, t := range tokens {
			valid := "\u2716"
			if t.ValidLink == "" {
				valid = "\u2713"
			}
			fmt.Printf("%- [1]*s  %s  %- [4]*s  %- [6]*s  % 10v  % 10v  % 10v %s\n",
				tokenLen, t.Token, valid, nameLen, t.Name, typeLen, t.Type,
				t.CreatedAt.Local().Format("2006-01-02"), t.UpdatedAt.Local().Format("2006-01-02"), t.ExpiredAt.Local().Format("2006-01-02"),
				t.ValidLink,
			)
		}
	},
}

func init() {
	tokenCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
