/*
Copyright © 2023 Vinyl Linux
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinyl-linux/mint/generator"
	"github.com/vinyl-linux/mint/parser"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate ./directory/of/mint-documents",
	Short: "Generate go code from mint documents",
	Long: `Generate parses a directory full of mint documents and
generates valid go code for Marshalling/ Unmarshalling types.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)

		if argCount != 1 {
			return fmt.Errorf("mint document directory missing")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		a, err := parser.ParseDir(args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			failer(errInvalid)
		}

		gen, _ := generator.New(a, &generator.GeneratorOptions{
			PackageName:             mustString(cmd.Flags().GetString("package")),
			Directory:               mustString(cmd.Flags().GetString("dest")),
			MakeDirectory:           mustBool(cmd.Flags().GetBool("mkdir")),
			CustomFunctionSkeletons: mustBool(cmd.Flags().GetBool("functions")),
			Clobber:                 mustBool(cmd.Flags().GetBool("clobber")),
		})

		err = gen.Generate()
		if err != nil {
			fmt.Printf("%v\n", err)
			failer(errInvalid)
		}
	},
}

func mustString(s string, err error) string {
	if err != nil {
		panic(err)
	}

	return s
}

func mustBool(b bool, err error) bool {
	if err != nil {
		panic(err)
	}

	return b
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	generateCmd.Flags().StringP("package", "p", "types", "Package to generate for")
	generateCmd.Flags().StringP("dest", "d", "types/", "Directory to generate code into")
	generateCmd.Flags().BoolP("mkdir", "m", true, "Create dest directory if not exist (note: if the destination directory exists then nothing happens)")
	generateCmd.Flags().BoolP("functions", "f", false, "Create skeleton functions for any custom validators and/or transforms found in documents")
	generateCmd.Flags().BoolP("clobber", "c", false, "Clobber any pre-generated validators and/or transforms (note: this replaces all custom code back to placeholders)")

}
