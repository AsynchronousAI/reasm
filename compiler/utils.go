package compiler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ReadDirective(directive string) []string {
	// Regex: match either "quoted strings" or sequences of non-comma, non-whitespace characters
	re := regexp.MustCompile(`"([^"]*)"|[^,\s]+`)
	matches := re.FindAllString(directive, -1)

	// Remove quotes and trim spaces
	for i, match := range matches {
		matches[i] = strings.TrimSpace(strings.Trim(match, `"`))
	}

	return matches
}
func AddEnd(w *OutputWriter) {
	if w.Depth == 0 {
		return
	}

	w.Depth--
	if w.Options.Comments {
		WriteIndentedString(w, "end -- %s (%s)\n", w.CurrentLabel, w.CurrentLabel)
	} else {
		WriteIndentedString(w, "end\n")
	}
}
func CompileRegister(w *OutputWriter, argument Argument) string {
	/* does it work as a integer (its a plain) */
	_, err := strconv.Atoi(argument.Source)
	if err == nil {
		return argument.Source
	}

	var compiled string = argument.Source

	if memoryAddress, ok := w.MemoryMap[argument.Source]; ok {
		if w.Options.Comments {
			compiled = fmt.Sprintf("%d --[[ %s ]]", memoryAddress, argument.Source)
		} else {
			compiled = fmt.Sprintf("%d", memoryAddress)
		}
	} else if isReg, regName := isRegister(argument.Source); isReg { /* it is a register! */
		regNumber := baseRegs[regName]
		if w.Options.Comments {
			compiled = fmt.Sprintf("registers[%d --[[ %s ]] ]", regNumber, regName)
		} else {
			compiled = fmt.Sprintf("registers[%d]", regNumber)
		}

		/** Offset */
		if argument.Offset != 0 {
			compiled = fmt.Sprintf("%s+%d", compiled, argument.Offset)
		}
	}

	/** Modifier */
	if argument.Modifier != "" {
		compiled = fmt.Sprintf("%s(%s)", argument.Modifier, compiled)
	}
	return compiled
}
func JumpTo(w *OutputWriter, label string, link bool) {
	address := FindLabelAddress(w, label)

	if address != -1 {
		WriteIndentedString(w, "do\n") // wrap with a do so luau does not complain if any code is after the continue
		w.Depth++
		if link {
			WriteIndentedString(w, "registers[2] = %d\n", w.MaxPC)
		}

		if w.Options.Comments {
			WriteIndentedString(w, "PC = %d -- %s\n", address, label)
		} else {
			WriteIndentedString(w, "PC = %d\n", address)
		}

		if w.Options.Trace {
			WriteIndentedString(w, "print('JUMP: ', PC)\n")
		}

		WriteIndentedString(w, "return true\n")
		w.Depth--
		WriteIndentedString(w, "end\n")
	} else {
		WriteIndentedString(w, "error(\"No bindings for functions '%s'\")\n", label)
	}
}
func CutAndLink(w *OutputWriter) {
	AddEnd(w)
	WriteIndentedString(w, "FUNCS[%d] = function() -- %s (extended) \n", w.MaxPC, w.CurrentLabel)
	w.Depth++
	w.MaxPC++
	w.CurrentLabel = IncrementFunctionName(w.CurrentLabel)
}
func FindInArray(array []string, target string) int {
	for i, item := range array {
		if item == target {
			return i
		}
	}
	return -1
}

func isCutoffInstruction(instruction AssemblyCommand) bool {
	return instruction.Type == Instruction && (instruction.Name == "call" || instruction.Name == "jal" || instruction.Name == "jalr")
}
func IncrementFunctionName(name string) string {
	re := regexp.MustCompile(`^(.*?)(?:_ext_(\d+))?$`)
	matches := re.FindStringSubmatch(name)

	if len(matches) == 0 {
		return name + "_ext_1"
	}

	base := matches[1]
	suffix := matches[2]

	if suffix == "" {
		return base + "_ext_1"
	}

	num, err := strconv.Atoi(suffix)
	if err != nil {
		return base + "_ext_1"
	}

	return fmt.Sprintf("%s_ext_%d", base, num+1)
}
func GetAllLabels(writer *OutputWriter) []string {
	labels := make([]string, 0)
	for _, command := range writer.Commands {
		if command.Type == Label {
			labels = append(labels, command.Name)
		}
	}
	return labels
}

func FindLabelAddress(writer *OutputWriter, target string) int {
	countedLabels := 0
	for _, label := range writer.Commands {
		if label.Type == Label || isCutoffInstruction(label) {
			countedLabels++ // our labels are Lua indexed starting at 1
		}
		if label.Type == Label && label.Name == target {
			return countedLabels
		}
	}
	return -1
}

func IsLabelEmpty(writer *OutputWriter, label string) bool {
	reachedLabel := false
	for _, command := range writer.Commands {
		if command.Type == Instruction && reachedLabel { /* found something in our bounds */
			return true
		}

		if command.Type == Label {
			if command.Name == label { /* it is our turn! */
				reachedLabel = true
			} else if reachedLabel { /* we passed another label, our time is up */
				return false
			}
		}
	}

	return false
}
