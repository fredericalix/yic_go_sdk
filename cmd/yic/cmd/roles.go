package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/youritcity/go-sdk/youritcity"

	"github.com/spf13/cobra"
)

var verbose bool

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List possible authorization roles",
	Run: func(cmd *cobra.Command, args []string) {
		conn := youritcity.NewConnection()
		roles, err := conn.GetRoles()
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't get roles:", err)
			os.Exit(1)
		}
		if len(roles) == 0 {
			fmt.Fprintln(os.Stderr, "no role received")
			os.Exit(1)
		}

		var keys []string
		for k := range roles {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Show list of roles
		if !verbose {
			for _, k := range keys {
				fmt.Println(k)
			}
			return
		}

		// Show list of roles with detailled view
		nameLen := 80
		rs := make(map[string]struct{})
		for k, role := range roles {
			if nameLen > len(k) {
				nameLen = len(k)
			}
			for rk := range role {
				rs[rk] = struct{}{}
			}
		}
		var rkeys []string
		for k := range rs {
			rkeys = append(rkeys, k)
		}
		sort.Strings(rkeys)

		// Header
		fmt.Print(" Roles  ")
		for _, k := range rkeys {
			fmt.Printf(" % 8s", k)
		}
		fmt.Println()

		for _, r := range keys {
			fmt.Printf("% 8s", r)
			for _, k := range rkeys {
				v := roles[r][k]
				switch v {
				case "":
					fmt.Print("         ")
				case "rw":
					fmt.Print("    rw   ")
				case "r":
					fmt.Print("    r    ")
				case "w":
					fmt.Print("     w   ")
				}
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(rolesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rolesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	rolesCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display more information on the roles")
}
