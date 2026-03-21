package compiler

import (
	"fmt"
	"strings"
)

type PendingDataType int8

const (
	PendingDataTypeNone    PendingDataType = 0
	PendingDataTypeString  PendingDataType = 1 /* a string generated via directive */
	PendingDataTypeNumeric PendingDataType = 2 /* a numeric value generated via .word for example */
)

type PendingData struct {
	Type PendingDataType
	Data string
}
type OutputWriter struct {
	Buffer                   []byte            /* the output */
	CurrentLabel             *AssemblyCommand  /* keep track of current label  */
	MemoryDevelopmentPointer int32             /* used when generating code that propagates memory with strings */
	PendingData              PendingData       /* used for remember data across instructions */
	Depth                    int               /* used for indentation */
	MaxPC                    int               /* used for counting PC which is hardcoded in */
	Commands                 []AssemblyCommand /* used to check lines in the future */
	IRNodes                  []*IRNode
	MemoryMap                map[string]int /* map static data keys to addresses */
	Options                  Options        /* user specified options */
	InlineCounter            int
	InstructionTotal         int
	InstructionProcessed     int
}

func newOutputWriter(options Options) *OutputWriter {
	return &OutputWriter{
		Buffer:                   []byte(""),
		MemoryDevelopmentPointer: 0,
		MaxPC:                    1,
		Options:                  options,
		MemoryMap:                make(map[string]int),
	}
}

func (w *OutputWriter) nextInlineTemp() string {
	name := fmt.Sprintf("__inline%d", w.InlineCounter)
	w.InlineCounter++
	return name
}

func WriteString(writer *OutputWriter, format string, args ...any) {
	writer.Buffer = append(writer.Buffer, fmt.Sprintf(format, args...)...)
}
func WriteIndentedString(writer *OutputWriter, format string, args ...any) {
	indent := strings.Repeat("\t", writer.Depth)
	writer.Buffer = append(writer.Buffer, indent+fmt.Sprintf(format, args...)...)
}

func (w *OutputWriter) updateProgress() {
	if w.InstructionTotal == 0 {
		return
	}
	fmt.Printf("\rInstruction %d/%d", w.InstructionProcessed, w.InstructionTotal)
}

func (w *OutputWriter) finishProgress() {
	if w.InstructionTotal == 0 {
		return
	}
	fmt.Print("\n")
}
