package compiler

import (
	"debug/elf"
	"fmt"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/arch/riscv64/riscv64asm"
)

func ParseFromElf(f *elf.File) []AssemblyCommand {
	/* validate file */
	if f.Class != elf.ELFCLASS32 {
		log.Error("ELF file is not 32-bit.")
		return []AssemblyCommand{}
	} else if f.Machine != elf.EM_RISCV {
		log.Error("ELF file is not RISC-V.")
		return []AssemblyCommand{}
	} else if f.Data != elf.ELFDATA2LSB {
		log.Error("ELF file is not little-endian.")
		return []AssemblyCommand{}
	} else if f.Type != elf.ET_EXEC {
		log.Error("ELF file is not an executable.")
		return []AssemblyCommand{}
	}

	/* look for bytecode */
	text := f.Section(".text")
	if text == nil {
		log.Error(".text section not found")
		return []AssemblyCommand{}
	}
	symtab, err := f.Symbols() // returns []elf.Symbol
	if err != nil {
		log.Println("No symbol table found:", err)
	}

	code, err := text.Data()
	if err != nil {
		log.Error(err)
		return []AssemblyCommand{}
	}

	sort.Slice(symtab, func(i, j int) bool {
		return symtab[i].Value < symtab[j].Value
	})

	instructions := make([]AssemblyCommand, 0)

	symIndex := 0
	for i := 0; i < len(code); i += 4 {
		inst, err := riscv64asm.Decode(code[i : i+4])
		if err != nil {
			fmt.Println("Decode error:", err)
			continue
		}

		addr := text.Addr + uint64(i)

		/* any symbols to be added before this */
		for symIndex < len(symtab) && symtab[symIndex].Value < addr {
			symIndex++
		}
		for symIndex < len(symtab) && symtab[symIndex].Value == addr {
			instructions = append(instructions, AssemblyCommand{
				Type:      Label,
				Name:      symtab[symIndex].Name,
				Arguments: nil,
			})
			symIndex++
		}

		addrLabel := fmt.Sprintf("%d", addr)
		if len(instructions) == 0 || instructions[len(instructions)-1].Name != addrLabel {
			instructions = append(instructions, AssemblyCommand{
				Type:      Label,
				Name:      addrLabel,
				Arguments: nil,
			})
		}

		/* read args */
		var argsStr []string
		for _, a := range inst.Args {
			if a == nil {
				continue
			}

			argsStr = append(argsStr, fmt.Sprintf("%s,", a))
		}

		args := parseArguments(argsStr)

		/* read instruction */
		instructionOp := strings.ToLower(inst.Op.String())

		/* append */
		instructions = append(instructions, AssemblyCommand{
			Type:      Instruction,
			Name:      instructionOp,
			Arguments: args,
		})
	}

	return instructions
}
