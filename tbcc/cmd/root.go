// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/backend"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/frontend"
	"github.com/thomasbollmeier/writing-a-c-compiler/tbcc/tacky"
	"os"
	"os/exec"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tbcc sourcefile",
	Short:   "A compiler for a simplified version of C",
	Long:    `TBCC is a compiler for a simplified version of C.`,
	Version: "0.5.1",
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
	stopAfterSemAnalysis  bool
	stopAfterIR           bool
	stopAfterCodegen      bool
	stopAfterCodeEmission bool
}

var (
	stopAfterLex          *bool = nil
	stopAfterParse        *bool = nil
	stopAfterSemAnalysis  *bool = nil
	stopAfterIR           *bool = nil
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
		if assemblyFile != "" && fileExists(assemblyFile) && !*stopAfterCodeEmission {
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
		*stopAfterSemAnalysis,
		*stopAfterIR,
		*stopAfterCodegen,
		*stopAfterCodeEmission,
	})

	if err != nil {
		return err
	}

	if assemblyFile == "" || *stopAfterCodeEmission {
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
	tokens, err := frontend.Tokenize(string(fileContent))
	if err != nil {
		return "", err
	}
	if options.stopAfterLex {
		return "", nil
	}

	// Run parser
	parser := frontend.NewParser(tokens)
	program, err := parser.ParseProgram()
	if err != nil {
		return "", err
	}
	if options.stopAfterParse {
		program.Accept(frontend.NewAstPrinter(4))
		return "", nil
	}

	// Semantic analysis
	nameCreator := frontend.NewNameCreator()
	program, err = frontend.AnalyzeSemantics(program, nameCreator)
	if err != nil {
		return "", err
	}
	if options.stopAfterSemAnalysis {
		program.Accept(frontend.NewAstPrinter(4))
		return "", nil
	}

	// Create TACKY
	emitter := tacky.NewTranslator(nameCreator)
	tackyProgram := emitter.Translate(program)
	if options.stopAfterIR {
		return "", nil
	}

	// assembly generation
	asmProgram := backend.NewTranslator().Translate(tackyProgram)
	if options.stopAfterCodegen {
		asmProgram.Accept(backend.NewAsmPrinter(4))
		return "", nil
	}

	// emit code
	assembly := backend.NewCodeGenerator().GenerateCode(*asmProgram)
	assemblyFile := stripSuffix(preProcessedFile) + ".s"
	err = os.WriteFile(assemblyFile, []byte(assembly), 0666)

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
	stopAfterSemAnalysis = rootCmd.PersistentFlags().Bool("validate", false, "stop after semantic analysis")
	stopAfterIR = rootCmd.PersistentFlags().Bool("tacky", false, "stop after tacky generation")
	stopAfterCodegen = rootCmd.PersistentFlags().Bool("codegen", false, "stop after codegen")
	stopAfterCodeEmission = rootCmd.PersistentFlags().BoolP("emission", "S", false, "stop after emission")
	rootCmd.MarkFlagsMutuallyExclusive("lex", "parse", "validate", "tacky", "codegen", "emission")
}
