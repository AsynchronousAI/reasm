package compiler

import (
	"debug/elf"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Options struct {
	Comments   bool
	Trace      bool
	Accurate   bool
	Memory     int
	Mode       string
	MainSymbol string
	Imports    []string
	LogIR      bool
}

func countInstructionTotal(commands []AssemblyCommand) int {
	total := 0
	for _, cmd := range commands {
		if cmd.Type == Instruction && cmd.Name != "" {
			total++
		}
	}
	return total
}

func Compile(executable *os.File, options Options) []byte {
	/* prepare */
	writer := newOutputWriter(options)

	elf, err := elf.NewFile(executable)
	if err != nil {
		assembly, _ := io.ReadAll(executable)
		assembly_str := string(assembly)
		lines := strings.Split(assembly_str, "\n")

		/* parse */
		for _, line := range lines {
			var command AssemblyCommand = Parse(writer, line)
			writer.Commands = append(writer.Commands, command)
		}
	} else {
		logrus.Warn(".elf support is experimental!")
		writer.Commands = ParseFromElf(elf)
	}

	writer.InstructionTotal = countInstructionTotal(writer.Commands)

	/* compilation */
	BeforeCompilation(writer)
	for _, command := range writer.Commands {
		CompileInstruction(writer, command)
		if command.Type == Instruction && command.Name != "" {
			writer.InstructionProcessed++
		}
		writer.updateProgress()
	}
	writer.finishProgress()
	if writer.Options.LogIR {
		dumpIRAsJSON(writer)
	}
	return AfterCompilation(writer)
}
