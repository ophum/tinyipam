/*
Copyright © 2024 Takahiro INAGAKI <inagaki0106@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/ophum/tinyipam/pkg/ipam"
	"github.com/ophum/tinyipam/pkg/ipam/file"
	"github.com/spf13/cobra"
)

var optsNewName = ""

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
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
		ip, err := s.AcquireIP(context.Background(), ipam.Name(optsNewName))
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s %s acquired\n", ip.ID, ip.IP().String())
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newCmd.Flags().StringVarP(&optsNewName, "name", "n", "", "name")
}
