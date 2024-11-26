/*
Copyright © 2024 Takahiro INAGAKI <inagaki0106@gmail.com>
*/
package cmd

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/ophum/tinyipam/pkg/ipam/file"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := file.New(dbFile)
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		ip, err := s.Get(ctx, args[0])
		if err != nil {
			panic(err)
		}

		w := tablewriter.NewWriter(os.Stdout)
		w.SetHeader([]string{"id", "ip", "name"})
		w.Append([]string{ip.ID, ip.IP().String(), ip.Name})
		w.Render()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
