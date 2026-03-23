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
	Buffer                   []byte           /* the output */
	CurrentLabel             *AssemblyCommand /* keep track of current label  */
	MemoryDevelopmentPointer int32            /* used when generating code that propagates memory with strings */
	PendingData              PendingData      /* used for remember data across instructions */
	Depth                    int              /* used for indentation */
	MaxPC                    int              /* used for counting PC which is hardcoded in */
	Commands                 []AssemblyCommand /* used to check lines in the future */
	IRNodes   []*IRNode
	MemoryMap map[string]int /* map static data keys to addresses */
	// LabelPC is the O(1) label-address cache built by BuildLabelCache.
	// Keyed by label name; value is the same sequential PC index that
	// FindLabelAddress previously computed via a linear scan.
	LabelPC map[string]int
	Options Options /* user specified options */
	InstructionTotal         int
	InstructionProcessed     int
}

// indentCache holds pre-built tab-indent strings for depths 0..maxCachedDepth.
// Compilation rarely exceeds a handful of levels, so 32 is a safe ceiling.
const maxCachedDepth = 32

var indentCache [maxCachedDepth]string

func init() {
	for i := range indentCache {
		indentCache[i] = strings.Repeat("\t", i)
	}
}

func cachedIndent(depth int) string {
	if depth < maxCachedDepth {
		return indentCache[depth]
	}
	return strings.Repeat("\t", depth)
}

func newOutputWriter(options Options) *OutputWriter {
	return &OutputWriter{
		// Pre-allocate a generous buffer to avoid repeated growth during compilation.
		// 256 KiB covers most small-to-medium programs; large ones will grow once or twice.
		Buffer:                   make([]byte, 0, 256*1024),
		MemoryDevelopmentPointer: 0,
		MaxPC:                    1,
		Options:                  options,
		MemoryMap:                make(map[string]int),
	}
}

// WriteString appends a formatted string to the buffer (no indentation).
func WriteString(writer *OutputWriter, format string, args ...any) {
	writer.Buffer = append(writer.Buffer, fmt.Sprintf(format, args...)...)
}

// WriteIndentedString prepends the current indentation level then appends the
// formatted string.  The indent string is looked up from a pre-built cache
// rather than allocated on every call.
func WriteIndentedString(writer *OutputWriter, format string, args ...any) {
	writer.Buffer = append(writer.Buffer, cachedIndent(writer.Depth)...)
	writer.Buffer = append(writer.Buffer, fmt.Sprintf(format, args...)...)
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