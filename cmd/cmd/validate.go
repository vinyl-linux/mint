/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vinyl-linux/mint/parser"
)

const errInvalid = 1

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate path/to/mint/documents/directory/",
	Short: "Validate the documents in a directory",
	Long: `Given a directory containing mint documents, read, parse
and validate the correctness of documents`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)

		if argCount != 1 {
			return fmt.Errorf("mint document directory missing")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		_, err := parser.ParseDir(args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(errInvalid)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
