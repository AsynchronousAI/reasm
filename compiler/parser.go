package compiler

import (
	"regexp"
	"strconv"
	"strings"
)

var offsetReg = regexp.MustCompile(`^(-?\d+)?\((\w+)\)$`)
var modifierReg = regexp.MustCompile(`^%(hi|lo)\((.+)\)$`)

type Argument struct {
	Offset   int8
	Register bool
	Source   string
	Modifier string /* "lo", "hi", or "" */
}
type CommandType uint8

const (
	Instruction CommandType = 0
	Label       CommandType = 1
	Directive   CommandType = 2
)

type AssemblyCommand struct {
	Type      CommandType
	Name      string
	Arguments []Argument
	Ignore    bool /* usually only used for labels, to exclude ones without instructions */
}

func Parse(writer *OutputWriter, command string) AssemblyCommand {
	// Remove comments
	if idx := strings.IndexAny(command, ";#"); idx != -1 {
		command = command[:idx]
	}

	cmd := strings.TrimSpace(command)

	// Empty line
	if cmd == "" {
		return AssemblyCommand{Type: Instruction}
	}

	// Label (ends with ':')
	if strings.HasSuffix(cmd, ":") {
		name := strings.TrimSuffix(cmd, ":")
		return AssemblyCommand{Type: Label, Name: name}
	}

	// Directive (.something ...)
	if strings.HasPrefix(cmd, ".") {
		name := cmd
		args := make([]Argument, 0)
		return AssemblyCommand{Type: Directive, Name: name, Arguments: args}
	}

	// Regular instruction (command)
	parts := strings.Fields(cmd)
	name := parts[0]
	args := parseArguments(parts[1:])

	return AssemblyCommand{Type: Instruction, Name: name, Arguments: args}
}

func parseArguments(parts []string) []Argument {
	if len(parts) == 0 {
		return nil
	}

	joined := strings.Join(parts, " ")
	pieces := strings.Split(joined, ",")

	args := []Argument{}
	for _, p := range pieces {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		arg := Argument{}

		// Check for %hi(...) or %lo(...) pattern
		if matches := modifierReg.FindStringSubmatch(p); matches != nil {
			arg.Modifier = matches[1] // "hi" or "lo"
			p = matches[2]            // content inside parentheses
		}

		arg.Source = p // store the raw source

		// Check for offset(register) pattern inside source
		isReg, reg := isRegister(p)
		if matches := offsetReg.FindStringSubmatch(p); matches != nil {
			arg.Register = true
			arg.Source = matches[2] // register name only
			if matches[1] != "" {
				if val, err := strconv.Atoi(matches[1]); err == nil {
					arg.Offset = int8(val)
				}
			} else {
				arg.Offset = 0
			}
		} else if isReg { // standalone register
			arg.Register = true
			arg.Offset = 0
			arg.Source = reg
		}

		args = append(args, arg)
	}
	return args
}

func isRegister(s string) (bool, string) {
	if baseRegs[s] > 0 {
		return true, s
	} else if abiToReg[s] != "" {
		return true, abiToReg[s]
	}
	return false, ""
}
