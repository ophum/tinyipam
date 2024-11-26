/*
Copyright © 2024 Takahiro INAGAKI <inagaki0106@gmail.com>
*/
package cmd

import (
	"context"
	"net"

	"github.com/ophum/tinyipam/pkg/ipam"
	"github.com/ophum/tinyipam/pkg/ipam/file"
	"github.com/spf13/cobra"
)

var updateOptions struct {
	name string
	ip   string
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
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

		if updateOptions.name != "" {
			ip.Name = updateOptions.name
		}
		if updateOptions.ip != "" {
			t := net.ParseIP(updateOptions.ip)
			ip.Addr = ipam.IPtoUint32(t)
		}
		if err := s.Update(ctx, ip); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().StringVarP(&updateOptions.name, "name", "n", "", "name")
	updateCmd.Flags().StringVarP(&updateOptions.ip, "ip", "i", "", "ip")
}
