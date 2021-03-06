package cmd

import (
	"fmt"
	"os"

	"github.com/fredericalix/yic_go_sdk/youritcity"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login your@email.com [role]",
	Short: "Login to yourITcity service",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t := "web"
		if len(args) >= 2 {
			t = args[1]
		}
		conn := youritcity.NewConnection()
		app, err := conn.Login(args[0], t)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't login:", err)
			os.Exit(1)
		}
		fmt.Printf("export %s=%s\n", EnvYourITcityToken, app.Token)
		fmt.Println("\nTo use other part if this cli run the commande above. Don't forget to validate this token by following the link sent you by email.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
