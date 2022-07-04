/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// colorvalidCmd represents the colorvalid command
var colorvalidCmd = &cobra.Command{
	Use:   "colorvalid",
	Short: "A brief description of your command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a color argument")
		}
		if IsValidColor(args[0]) {
			return nil
		}

		return fmt.Errorf("invalid color specified: %s", args[0])
	},

	Run: func(cmd *cobra.Command, args []string) {
		if cmd.("t") {

		}
		fmt.Println("colorvalid called, color is:", args[0], "t:")
	},
}

var Ttt bool

func init() {
	rootCmd.AddCommand(colorvalidCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// colorvalidCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	colorvalidCmd.Flags().BoolVarP(&Ttt, "toggle", "t", false, "Help message for toggle")
}

func IsValidColor(color string) bool {
	if color == "red" {
		return true
	}

	return false
}
