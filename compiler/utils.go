package compiler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var symbolWithOffset = regexp.MustCompile(`^([.$A-Za-z_][.$A-Za-z0-9_]*)([+-]\d+)$`)
var incrementFunctionNameRe = regexp.MustCompile(`^(.*?)(?:_ext_(\d+))?$`)

func ReadDirective(directive string) []string {
	result := make([]string, 0, 4)
	i := 0
	for i < len(directive) {
		for i < len(directive) && (directive[i] == ' ' || directive[i] == '\t' || directive[i] == ',') {
			i++
		}
		if i >= len(directive) {
			break
		}

		if directive[i] == '"' {
			i++ // consume opening quote
			var sb strings.Builder
			escaped := false
			for i < len(directive) {
				ch := directive[i]
				i++
				if escaped {
					sb.WriteByte(ch)
					escaped = false
					continue
				}
				if ch == '\\' {
					sb.WriteByte(ch)
					escaped = true
					continue
				}
				if ch == '"' {
					break
				}
				sb.WriteByte(ch)
			}
			result = append(result, sb.String())
			continue
		}

		start := i
		for i < len(directive) && directive[i] != ',' && directive[i] != ' ' && directive[i] != '\t' {
			i++
		}
		token := strings.TrimSpace(directive[start:i])
		if token != "" {
			result = append(result, token)
		}
	}
	return result
}

func UnescapeDirectiveString(raw string) (string, error) {
	return strconv.Unquote(`"` + raw + `"`)
}

func luauLongBracketLiteral(content string) string {
	delimiter := ""
	for strings.Contains(content, "]"+delimiter+"]") {
		delimiter += "="
	}
	return fmt.Sprintf("[%s[%s]%s]", delimiter, content, delimiter)
}

func luauStringExpression(content string) string {
	if content == "" {
		return `[[]]`
	}

	parts := make([]string, 0, 4)
	bytes := []byte(content)
	for i := 0; i < len(bytes); {
		if bytes[i] == 0 {
			run := 1
			for i+run < len(bytes) && bytes[i+run] == 0 {
				run++
			}
			if run == 1 {
				parts = append(parts, `"\0"`)
			} else {
				parts = append(parts, fmt.Sprintf(`string.rep("\0", %d)`, run))
			}
			i += run
			continue
		}

		start := i
		for i < len(bytes) && bytes[i] != 0 {
			i++
		}
		parts = append(parts, luauLongBracketLiteral(string(bytes[start:i])))
	}
	return strings.Join(parts, " .. ")
}

func AddEnd(w *OutputWriter) {
	if w.Depth == 0 {
		return
	}

	WriteIndentedString(w, "return false\n")
	w.Depth--
	if w.Options.Comments {
		WriteIndentedString(w, "end -- %s (%s)\n", w.CurrentLabel.Name, w.CurrentLabel.Name)
	} else {
		WriteIndentedString(w, "end\n")
	}
}

func CompileRegister(w *OutputWriter, argument Argument) string {
	/* does it work as a integer (its a plain) */
	_, err := strconv.ParseInt(argument.Source, 0, 64)
	if err == nil {
		return argument.Source
	}

	var compiled string = argument.Source

	if memoryAddress, ok := resolveSymbolAddress(w, argument.Source); ok {
		if argument.Modifier != "" {
			/* resolve to plain number for use inside modifier call, add comment outside */
			inner := fmt.Sprintf("%d", memoryAddress)
			if literal, ok := resolveModifierLiteral(argument.Modifier, memoryAddress); ok {
				compiled = fmt.Sprintf("%d", literal)
			} else {
				compiled = applyModifierExpression(argument.Modifier, inner)
			}
			if argument.BaseRegister != "" {
				baseNum := baseRegs[argument.BaseRegister]
				if w.Options.Comments {
					compiled = fmt.Sprintf("%s + --[[ %s ]] %s", compiled, argument.BaseRegister, regVarName(baseNum))
				} else {
					compiled = fmt.Sprintf("%s + %s", compiled, regVarName(baseNum))
				}
			}
			if w.Options.Comments {
				compiled = fmt.Sprintf("--[[ %s ]] %s", argument.Source, compiled)
			}
			return compiled
		}
		if w.Options.Comments {
			compiled = fmt.Sprintf("--[[ %s ]] %d", argument.Source, memoryAddress)
		} else {
			compiled = fmt.Sprintf("%d", memoryAddress)
		}
	} else if isReg, regName := isRegister(argument.Source); isReg { /* it is a register! */
		regNumber := baseRegs[regName]
		if w.Options.Comments {
			compiled = fmt.Sprintf("--[[ %s ]] %s", regName, regVarName(regNumber))
		} else {
			compiled = regVarName(regNumber)
		}

		/** Offset */
		if argument.Offset != 0 {
			compiled = fmt.Sprintf("%s+%d", compiled, argument.Offset)
		}
	} else {
		logrus.Errorf("Could not resolve as symbol or register: %s\n", argument.Source)
		compiled = "0"
	}

	/** Modifier (source was not in MemoryMap — leave as symbol name) */
	if argument.Modifier != "" {
		compiled = applyModifierExpression(argument.Modifier, compiled)
	}

	/** Base Register for %lo(sym)(reg) */
	if argument.BaseRegister != "" {
		baseNum := baseRegs[argument.BaseRegister]
		var baseCompiled string
		if w.Options.Comments {
			baseCompiled = fmt.Sprintf("--[[ %s ]] %s", argument.BaseRegister, regVarName(baseNum))
		} else {
			baseCompiled = regVarName(baseNum)
		}
		compiled = fmt.Sprintf("%s + %s", compiled, baseCompiled)
	}
	return compiled
}

func resolveModifierLiteral(modifier string, address int) (int, bool) {
	value := uint32(address)
	switch modifier {
	case "hi":
		return int((value + 0x800) >> 12), true
	case "lo":
		raw := int32(value & 0xFFF)
		if raw >= 0x800 {
			raw -= 0x1000
		}
		return int(raw), true
	default:
		return 0, false
	}
}

func applyModifierExpression(modifier, expr string) string {
	switch modifier {
	case "hi":
		return fmt.Sprintf("bit32.rshift(bit32.band((%s) + 0x800, 0xFFFFFFFF), 12)", expr)
	case "lo":
		return fmt.Sprintf("bit32.band((%s) + 0x800, 0xFFF) - 0x800", expr)
	default:
		return fmt.Sprintf("%s(%s)", modifier, expr)
	}
}

func resolveSymbolAddress(w *OutputWriter, symbol string) (int, bool) {
	if memoryAddress, ok := w.MemoryMap[symbol]; ok {
		return memoryAddress, true
	}

	matches := symbolWithOffset.FindStringSubmatch(symbol)
	if matches == nil {
		return 0, false
	}

	offset, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, false
	}

	baseAddress, ok := w.MemoryMap[matches[1]]
	if !ok {
		return 0, false
	}

	return baseAddress + offset, true
}

func JumpTo(w *OutputWriter, label string, link bool) {
	address := FindLabelAddress(w, label)

	if address != -1 {
		WriteIndentedString(w, "do\n") // wrap with a do so luau does not complain if any code is after the continue
		w.Depth++
		if link {
			WriteIndentedString(w, "r2 = %d\n", w.MaxPC)
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
	WriteIndentedString(w, "FUNCS[%d] = function(): boolean -- %s (extended) \n", w.MaxPC, w.CurrentLabel.Name)
	w.Depth++
	w.MaxPC++
	w.CurrentLabel.Name = IncrementFunctionName(w.CurrentLabel.Name)
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
	return instruction.Type == Instruction && (instruction.Name == "call" || instruction.Name == "tail" || instruction.Name == "jal" || instruction.Name == "jalr")
}

func IncrementFunctionName(name string) string {
	matches := incrementFunctionNameRe.FindStringSubmatch(name)

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
	labels := make([]string, 0, len(writer.LabelPC))
	for _, command := range writer.Commands {
		if command.Type == Label && !command.Ignore {
			labels = append(labels, command.Name)
		}
	}
	return labels
}

func BuildLabelCache(writer *OutputWriter) {
	writer.LabelPC = make(map[string]int, len(writer.Commands)/4+8)

	counted := 0
	for _, cmd := range writer.Commands {
		if (cmd.Type == Label && !cmd.Ignore) || isCutoffInstruction(cmd) {
			counted++
		}
		if cmd.Type == Label && !cmd.Ignore {
			writer.LabelPC[cmd.Name] = counted
		}
	}
}

func FindLabelAddress(writer *OutputWriter, target string) int {
	if writer.LabelPC != nil {
		if addr, ok := writer.LabelPC[target]; ok {
			return addr
		}
		return -1
	}
	counted := 0
	for _, label := range writer.Commands {
		if (label.Type == Label && !label.Ignore) || isCutoffInstruction(label) {
			counted++
		}
		if label.Type == Label && !label.Ignore && label.Name == target {
			return counted
		}
	}
	return -1
}

// ---------------------------------------------------------------------------
// IR integration helpers
// ---------------------------------------------------------------------------
func irArgExpr(w *OutputWriter, arg Argument) *IRNode {
	compiled := CompileRegister(w, arg)
	return &IRNode{Kind: IRExprLit, Op: compiled}
}
func JumpToIR(w *OutputWriter, cond *IRNode, label string) {
	address := FindLabelAddress(w, label)
	if address == -1 {
		Emit(w, IRStmtError("No bindings for function '"+label+"'"))
		return
	}
	var pcStmt *IRNode
	if w.Options.Comments {
		pcStmt = IRStmtAssignComment(IRSymbol(SYM_PC), IRLit(address), label)
	} else {
		pcStmt = IRStmtSetPC(IRLit(address))
	}
	body := []*IRNode{pcStmt}
	if w.Options.Trace {
		body = append(body, IRStmtCall("print", IRLit(`"JUMP: "`), IRSymbol(SYM_PC)))
	}
	body = append(body, IRStmtReturn(true))
	Emit(w, IRStmtIf(cond, body, nil))
}
