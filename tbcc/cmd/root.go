// Package cmd /*
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tbcc sourcefile",
	Short:   "A compiler for a simplified version of C",
	Long:    `TBCC is a compiler for a simplified version of C.`,
	Version: "0.1.0",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args[0])
		fmt.Println(*stopAfterLex)
		fmt.Println(*stopAfterParse)
		fmt.Println(*stopAfterCodegen)
	},
}

var (
	stopAfterLex     *bool = nil
	stopAfterParse   *bool = nil
	stopAfterCodegen *bool = nil
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	stopAfterLex = rootCmd.PersistentFlags().Bool("lex", false, "stop after lexer")
	stopAfterParse = rootCmd.PersistentFlags().Bool("parse", false, "stop after parser")
	stopAfterCodegen = rootCmd.PersistentFlags().Bool("codegen", false, "stop after codegen")
	rootCmd.MarkFlagsMutuallyExclusive("lex", "parse", "codegen")
}
