package main

import (
	"os"
	"slices"
	"strings"

	"github.com/AsynchronousAI/reasm/compiler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var enableComments bool
	var enableTrace bool
	var mode string
	var outputFile string
	var mainSymbol string
	var importSymbols []string // <- multiple imports

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true, // force color output
		FullTimestamp: true, // show timestamps
	})

	log.SetLevel(log.DebugLevel)

	var rootCmd = &cobra.Command{
		Use:   "reasm [input] [output]",
		Short: "Compile RISC-V Assembly into Luau",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, inputFiles []string) error {
			validModes := []string{"module", "main", "bench"}
			modeLower := strings.ToLower(mode)
			if !slices.Contains(validModes, modeLower) {
				log.Error("invalid mode. Valid modes are: module, main, bench")
				return nil
			}

			/* read input file */
			if len(inputFiles) > 1 {
				log.Error("Only one input file is supported at the moment, if you want to compile multiple files link before hand using an ELF file.")
				return nil
			}

			file, err := os.Open(inputFiles[0])
			if err != nil {
				log.Errorf("failed to read input file: %v", err)
				return nil
			}
			defer file.Close()

			/* compile with options */
			processed := compiler.Compile(file, compiler.Options{
				Comments:   enableComments,
				Trace:      enableTrace,
				Mode:       modeLower,
				MainSymbol: mainSymbol,
				Imports:    importSymbols,
			})

			/* write output file */
			err = os.WriteFile(outputFile, processed, 0644)
			if err != nil {
				log.Errorf("failed to write output file: %v", err)
				return nil
			}

			return nil
		},
	}

	// Flags
	rootCmd.Flags().BoolVar(&enableComments, "comments", false, "Include debug comments in the output")
	rootCmd.Flags().BoolVar(&enableTrace, "trace", false, "Prints out a trace of the PC")
	rootCmd.Flags().StringVar(&mode, "mode", "main", "Mode to compile as: module, main, or bench")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "The output luau file.")
	rootCmd.Flags().StringVarP(&mainSymbol, "symbol", "e", "main", "The main symbol to start automatically.")
	rootCmd.Flags().StringArrayVarP(&importSymbols, "import", "I", []string{}, "Import symbol(s), can be repeated (example: -Imath -Ios)")
	rootCmd.MarkFlagRequired("o")

	rootCmd.Execute()
}
