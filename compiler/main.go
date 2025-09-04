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
	Mode       string
	MainSymbol string
	Imports    []string
}

func Compile(executable *os.File, options Options) []byte {
	/* prepare */
	var writer = &OutputWriter{
		Buffer:                   []byte(""),
		MemoryDevelopmentPointer: 0,
		MaxPC:                    1,
		Options:                  options,
		MemoryMap:                make(map[string]int),
	}

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

	/* compilation */
	BeforeCompilation(writer)
	for _, command := range writer.Commands {
		CompileInstruction(writer, command)
	}
	return AfterCompilation(writer)
}
