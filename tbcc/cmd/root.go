// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
	"os"
	"os/exec"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tbcc sourcefile",
	Short:   "A compiler for a simplified version of C",
	Long:    `TBCC is a compiler for a simplified version of C.`,
	Version: "0.1.1",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := run(args)
		if err != nil {
			os.Exit(1)
		}
	},
}

type Options struct {
	stopAfterLex          bool
	stopAfterParse        bool
	stopAfterCodegen      bool
	stopAfterCodeEmission bool
}

var (
	stopAfterLex          *bool = nil
	stopAfterParse        *bool = nil
	stopAfterCodegen      *bool = nil
	stopAfterCodeEmission *bool = nil
)

func run(args []string) error {

	var preProcessedFile string
	var assemblyFile string

	defer func() {
		if preProcessedFile != "" && fileExists(preProcessedFile) {
			_ = os.Remove(preProcessedFile)
		}
		if assemblyFile != "" && fileExists(assemblyFile) {
			_ = os.Remove(assemblyFile)
		}
	}()

	preProcessedFile, err := preProcess(args[0])
	if err != nil {
		return err
	}

	assemblyFile, err = compile(preProcessedFile, Options{
		*stopAfterLex,
		*stopAfterParse,
		*stopAfterCodegen,
		*stopAfterCodeEmission,
	})

	if err != nil {
		return err
	} else if assemblyFile == "" {
		return nil
	}

	_, err = assemble(assemblyFile)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func assemble(assemblyFile string) (string, error) {
	execFile := stripSuffix(assemblyFile)
	cmd := exec.Command("gcc", assemblyFile, "-o", execFile)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return execFile, nil
}

func compile(preProcessedFile string, options Options) (string, error) {
	fileContent, err := os.ReadFile(preProcessedFile)
	if err != nil {
		return "", err
	}

	// Run lexer
	_, err = frontend.Tokenize(string(fileContent))
	if err != nil {
		return "", err
	}
	if options.stopAfterLex {
		return "", nil
	}

	assemblyFile := stripSuffix(preProcessedFile) + ".s"
	cmd := exec.Command("gcc", "-S", preProcessedFile, "-o", assemblyFile)
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return assemblyFile, nil
}

func preProcess(sourceFile string) (string, error) {
	preProcessedFile := stripSuffix(sourceFile) + ".i"
	cmd := exec.Command("gcc", "-E", "-P", sourceFile, "-o", preProcessedFile)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return preProcessedFile, nil
}

func stripSuffix(filename string) string {
	parts := strings.Split(filename, ".")
	length := len(parts)
	parts = parts[:length-1]
	return strings.Join(parts, ".")
}

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
	stopAfterCodeEmission = rootCmd.PersistentFlags().BoolP("emission", "S", false, "stop after emission")
	rootCmd.MarkFlagsMutuallyExclusive("lex", "parse", "codegen", "emission")
}
