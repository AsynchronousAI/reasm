package compiler

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

//go:embed templates/boilerplate.luau
var luau_boilerplate string

var instructions = map[string]func(*OutputWriter, AssemblyCommand){
	/* bit shifts */
	"sll":  sll,
	"srl":  srl,
	"slli": sll,
	"srli": srl,
	"sra":  sra,
	"srai": sra,

	/* bit operations */
	"and": and,
	"xor": xor,
	"or":  or,
	"not": not,

	/** immediate */
	"andi": and,
	"xori": xor,
	"ori":  or,

	/* memory */
	/** save */
	"sb": sb,
	"sh": sh,
	"sw": sw,

	/** load */
	"lb": lb,
	"lh": lh,
	"lw": lw,

	/*** variants */
	"li":  li,
	"lui": lui,
	"lbu": lbu,
	"lhu": lhu,

	/* math */
	"add":  add,
	"addi": add,
	"sub":  sub,
	"subi": sub,
	"neg":  neg,

	/* M extension */
	"div": div,
	"mul": mul,
	"rem": rem,

	/*** descendants */
	"remu":  rem,
	"mulh":  mulh,
	"mulhu": mulh,
	"mulsu": mulh,
	"mulu":  mulh,
	"divu":  div,

	/* branching */
	"bne":  bne,
	"blt":  blt,
	"bltu": blt,
	"bge":  bge,
	"beq":  beq,
	"bgeu": bge,
	"bgt":  bgt,
	"bgtu": bgt,
	"ble":  ble,
	"bleu": ble,

	/** zero descendants */
	"bnez": bnez,
	"beqz": beqz,
	"bltz": bltz,
	"bgtz": bgtz,
	"blez": blez,
	"bgez": bgez,

	/* jump */
	"j":    jump,
	"jalr": jalr,
	"jr":   jr,
	"jal":  jal,

	/* os */
	"ecall":  ecall,
	"ebreak": ebreak,
	"fence":  fence,

	/* set less/greator then */
	"slt":   slt,
	"sltu":  slt,
	"sltiu": slt,
	"slti":  slt,
	"seqz":  seqz,
	"snez":  snez,
	"sgtz":  sgtz,
	"sltz":  sltz,
	"sgt":   sgt,
	"sgtu":  sgt,

	/* F extension */
	"fclass.s": fclass,
	"fclass.d": fclass,

	/** Arithmetic */
	"fadd.s": fadd,
	"fsub.s": fsub,
	"fdiv.s": fdiv,
	"fmul.s": fmul,
	"fadd.d": fadd,
	"fsub.d": fsub,
	"fdiv.d": fdiv,
	"fmul.d": fmul,
	"fneg.s": fneg,
	"fneg.d": fneg,

	/** More advanced */
	"fabs.s":  fabs,
	"fabs.d":  fabs,
	"fsqrt.s": fsqrt,
	"fmin.s":  fmin,
	"fmax.s":  fmax,
	"fsqrt.d": fsqrt,
	"fmin.d":  fmin,
	"fmax.d":  fmax,

	/** Memory */
	"flw": flw,
	"fsw": fsw,
	"fld": fld,
	"fsd": fsd,

	/** Sign */
	"fsgnj.s":  fsgnj,
	"fsgnjn.s": fsgnjn,
	"fsgnjx.s": fsgnjx,
	"fsgnj.d":  fsgnj,
	"fsgnjn.d": fsgnjn,
	"fsgnjx.d": fsgnjx,

	/** Comparators */
	"feq.s": feq,
	"flt.s": flt,
	"fle.s": fle,
	"fgt.s": fgt,
	"fge.s": fge,
	"feq.d": feq,
	"flt.d": flt,
	"fle.d": fle,
	"fgt.d": fgt,
	"fge.d": fge,

	/** Fused */
	"fmadd.s":  fmadd,
	"fmsub.s":  fmsub,
	"fnmadd.s": fnmadd,
	"fnmsub.s": fnmsub,

	"fmadd.d":  fmadd,
	"fmsub.d":  fmsub,
	"fnmadd.d": fnmadd,
	"fnmsub.d": fnmsub,

	/** Conversion */
	"fmv.d": move,
	"fmv.s": move,

	"fmv.s.x":   fmv_w_x,
	"fmv.w.x":   fmv_w_x,
	"fmv.x.w":   fmv_x_w,
	"frflags":   frflags,
	"fsflags":   fsflags,
	"fcvt.w.s":  fcvt_w_s,
	"fcvt.wu.s": fcvt_w_s,
	"fcvt.s.w":  fcvt_s_w,
	"fcvt.s.wu": fcvt_s_wu,
	"fcvt.d.s":  fcvt_d_s,
	"fcvt.s.d":  fcvt_d_s,
	"fcvt.w.d":  fcvt_w_d,
	"fcvt.d.w":  fcvt_d_w,
	"fcvt.d.wu": fcvt_d_wu,

	/* Abstraction */
	"auipc": auipc,
	"ret":   ret,
	"call":  call,
	"tail":  tail,
	"mv":    move,
	"nop":   nop,

	/* Zbb Extension - Basic Bit Manipulation */
	/** Count operations */
	"clz":  clz,
	"ctz":  ctz,
	"cpop": cpop,

	/** Min/Max */
	"min":  min,
	"minu": minu,
	"max":  max,
	"maxu": maxu,

	/** Sign/Zero Extension */
	"sext.b": sext_b,
	"sext.h": sext_h,
	"zext.h": zext_h,

	/** Logical with Negate */
	"andn": andn,
	"orn":  orn,
	"xnor": xnor,

	/** Rotation */
	"rol":  rol,
	"ror":  ror,
	"rori": rori,

	/** Byte Operations */
	"orc.b": orc_b,
	"rev8":  rev8,

	/** Packing */
	"pack":  pack,
	"packh": packh,

	/* Zbs Extension - Single-Bit Instructions */
	/** Bit Set */
	"bset":  bset,
	"bseti": bseti,

	/** Bit Clear */
	"bclr":  bclr,
	"bclri": bclri,

	/** Bit Invert */
	"binv":  binv,
	"binvi": binvi,

	/** Bit Extract */
	"bext":  bext,
	"bexti": bexti,
}
var directives = map[string]func(*OutputWriter, []string){
	".asciz":   asciz,
	".string":  asciz,
	".base64":  base64data,
	".quad":    quad,
	".word":    word,
	".byte":    byte_,
	".half":    half,
	".zero":    zero,
	".align":   align,
	".p2align": align,
	".set":     set,
}

/* main */
func CompileInstruction(writer *OutputWriter, command AssemblyCommand) {
	switch command.Type {
	case Instruction:
		if command.Name == "" {
			break
		}

		if cmdFunc, ok := instructions[command.Name]; ok {
			if writer.Options.Comments {
				WriteIndentedString(writer, "-- %s (%v)\n", command.Name, command.Arguments)
			}

			cmdFunc(writer, command)
		} else {
			log.Warn("unknown instruction: " + command.Name)
		}
	case Label:
		label(writer, command)
	}
}
func normalizeSection(name string) string {
	name = strings.Split(name, ",")[0]
	name = strings.Trim(name, "\"")
	if strings.HasPrefix(name, ".rodata") {
		return ".rodata"
	}
	if strings.HasPrefix(name, ".sdata") {
		return ".sdata"
	}
	if strings.HasPrefix(name, ".data") {
		return ".data"
	}
	if strings.HasPrefix(name, ".sbss") {
		return ".sbss"
	}
	if strings.HasPrefix(name, ".bss") {
		return ".bss"
	}
	if strings.HasPrefix(name, ".text") {
		return ".text"
	}
	return ".rodata"
}

func BeforeCompilation(writer *OutputWriter) {
	/* Pass 1: Collect all labels into MemoryMap */
	tempMaxPC := 1
	var pendingLabels []string

	memSize := int32(writer.Options.Memory)
	sectionPointers := map[string]int32{
		".rodata": 1024,
		".data":   memSize / 8,
		".sdata":  2 * memSize / 8,
		".bss":    3 * memSize / 8,
		".sbss":   4 * memSize / 8,
		".text":   5 * memSize / 8,
	}
	currentSection := ".text"

	// Provide dummy allocations for common newlib/C externs
	writer.MemoryMap["_impure_ptr"] = 16
	writer.MemoryMap["_ctype_"] = 32

	for i := range writer.Commands {
		if writer.Commands[i].Type == Label {
			labelName := writer.Commands[i].Name
			pendingLabels = append(pendingLabels, labelName)
			writer.MemoryMap[labelName] = tempMaxPC // default to current MaxPC
		} else if writer.Commands[i].Type == Instruction && writer.Commands[i].Name != "" {
			for _, l := range pendingLabels {
				writer.MemoryMap[l] = tempMaxPC
			}
			pendingLabels = nil
			tempMaxPC++
		} else if writer.Commands[i].Type == Directive {
			attributeComponents := ReadDirective(writer.Commands[i].Name)
			attributeName := attributeComponents[0]

			if attributeName == ".section" || attributeName == ".text" || attributeName == ".data" || attributeName == ".bss" || attributeName == ".sdata" || attributeName == ".sbss" {
				sec := attributeName
				if attributeName == ".section" && len(attributeComponents) > 1 {
					sec = attributeComponents[1]
				}
				currentSection = normalizeSection(sec)
				continue
			}

			tempMemPtr := sectionPointers[currentSection]

			for _, l := range pendingLabels {
				writer.MemoryMap[l] = int(tempMemPtr)
			}
			pendingLabels = nil

			// Some directives increment MemoryDevelopmentPointer
			if attributeName == ".asciz" || attributeName == ".string" {
				data, _ := UnescapeDirectiveString(attributeComponents[1])
				tempMemPtr += int32(len(data) + 1)
			} else if attributeName == ".word" {
				tempMemPtr += 4
			} else if attributeName == ".half" {
				tempMemPtr += 2
			} else if attributeName == ".byte" {
				tempMemPtr += 1
			} else if attributeName == ".quad" {
				tempMemPtr += 8
			} else if attributeName == ".zero" {
				size, _ := strconv.Atoi(attributeComponents[1])
				tempMemPtr += int32(size)
			} else if attributeName == ".align" || attributeName == ".p2align" {
				pow, _ := strconv.Atoi(attributeComponents[1])
				alignSize := int32(1 << pow)
				rem := tempMemPtr % alignSize
				if rem != 0 {
					tempMemPtr += alignSize - rem
				}
			} else if attributeName == ".set" {
				symName := attributeComponents[1]
				if len(attributeComponents) >= 3 {
					if attributeComponents[2] == "." && len(attributeComponents) >= 5 {
						offset, _ := strconv.Atoi(attributeComponents[4])
						if attributeComponents[3] == "+" {
							writer.MemoryMap[symName] = int(tempMemPtr) + offset
						} else if attributeComponents[3] == "-" {
							writer.MemoryMap[symName] = int(tempMemPtr) - offset
						}
					} else {
						val, err := strconv.Atoi(attributeComponents[2])
						if err == nil {
							writer.MemoryMap[symName] = val
						} else {
							if existingAddr, ok := writer.MemoryMap[attributeComponents[2]]; ok {
								writer.MemoryMap[symName] = existingAddr
							}
						}
					}
				}
			}
			
			sectionPointers[currentSection] = tempMemPtr
		}
	}

	/* Pass 2: Actually process directives and write init() function */
	WriteIndentedString(writer, "function init(): ()\n")
	writer.Depth++
	WriteIndentedString(writer, "reset_registers()\n")

	sectionPointers = map[string]int32{
		".rodata": 1024,
		".data":   memSize / 8,
		".sdata":  2 * memSize / 8,
		".bss":    3 * memSize / 8,
		".sbss":   4 * memSize / 8,
		".text":   5 * memSize / 8,
	}
	currentSection = ".text"
	writer.MemoryDevelopmentPointer = sectionPointers[currentSection]

	for i := range writer.Commands {
		if writer.Commands[i].Type == Label {
			writer.CurrentLabel = &writer.Commands[i]
			writer.Commands[i].Ignore = true
			writer.PendingData = PendingData{}
		}
		if writer.Commands[i].Type == Instruction && writer.Commands[i].Name != "" {
			if writer.CurrentLabel != nil {
				writer.CurrentLabel.Ignore = false
			}
		}
		if writer.Commands[i].Type == Directive {
			attributeComponents := ReadDirective(writer.Commands[i].Name)
			attributeName := attributeComponents[0]

			if attributeName == ".section" || attributeName == ".text" || attributeName == ".data" || attributeName == ".bss" || attributeName == ".sdata" || attributeName == ".sbss" {
				sectionPointers[currentSection] = writer.MemoryDevelopmentPointer
				sec := attributeName
				if attributeName == ".section" && len(attributeComponents) > 1 {
					sec = attributeComponents[1]
				}
				currentSection = normalizeSection(sec)
				writer.MemoryDevelopmentPointer = sectionPointers[currentSection]
				continue
			}

			if _, ok := directives[attributeName]; ok {
				directives[attributeName](writer, attributeComponents)
			} else if writer.Options.Comments {
				WriteIndentedString(writer, "-- ASM DIRECTIVE: %s\n", writer.Commands[i].Name)
			}
		}
	}
	sectionPointers[currentSection] = writer.MemoryDevelopmentPointer

	/* reset MaxPC for second pass */
	writer.MaxPC = 1

	/* finish loading directives */
	WriteIndentedString(writer, "PC = %d\n", FindLabelAddress(writer, writer.Options.MainSymbol))
	WriteIndentedString(writer, "r3 = buffer.len(memory) - 1024 -- start at the end of memory minus some padding\n")
	WriteIndentedString(writer, "if r3 <= 0 then error(\"Not enough memory\") end\n")
	writer.Depth--
	WriteIndentedString(writer, "end\n")

	log.Infof("Section ends: rodata=%d, data=%d, sdata=%d, bss=%d, sbss=%d, text=%d",
		sectionPointers[".rodata"], sectionPointers[".data"], sectionPointers[".sdata"],
		sectionPointers[".bss"], sectionPointers[".sbss"], sectionPointers[".text"])
}

func AfterCompilation(writer *OutputWriter) []byte {
	AddEnd(writer) // end the current label, if active

	// check if invalid PC, then break
	WriteIndentedString(writer, "function start(startPosition: number): ()\n")
	writer.Depth++
	WriteIndentedString(writer, "PC = startPosition\n")
	WriteIndentedString(writer, "while FUNCS[PC] do\n")
	writer.Depth++
	WriteIndentedString(writer, "if not FUNCS[PC]() then\n")
	writer.Depth++
	WriteIndentedString(writer, "PC += 1\n")
	writer.Depth--
	WriteIndentedString(writer, "end\n")
	if writer.Options.Trace {
		WriteIndentedString(writer, "print(\"FALL THROUGH:\", PC)\n")
	}
	writer.Depth--
	WriteIndentedString(writer, "end\n")
	WriteIndentedString(writer, "flush_stdout()\n")
	writer.Depth--
	WriteIndentedString(writer, "end\n")

	// final code based on mode
	main := FindLabelAddress(writer, writer.Options.MainSymbol)
	if writer.Options.Mode == "main" {
		WriteString(writer, "init()\nstart(%d)\n", main)
	} else if writer.Options.Mode == "module" {
		WriteString(writer, "init()\n")
		WriteString(writer, `return setmetatable({
	init = init,
	memory = memory,
	functions = functions,
	files = files,
	util = {
		get_args = get_args,
		push_args = push_args,
		get_f_args = get_f_args,
		push_f_args = push_f_args,
		read_string = read_string,
		format_string = format_string,
		malloc = malloc,
		two_words_to_double = two_words_to_double,
		reset_registers = reset_registers,
	},


	exports = {
`)

		for _, label := range GetAllLabels(writer) {
			WriteString(writer, "\t\t[\"%s\"] = function() start(%d) end,\n", label, FindLabelAddress(writer, label))
		}

		WriteString(writer, "\t}\n}, {__call = function() init(); start(%d) end})", main)
	} else if writer.Options.Mode == "bench" {
		WriteString(writer, `
return {
    Name = "RISCV File",

    BeforeEach = init,

    Functions = {
        ["main"] = function() start(%d) end,
    }
}`, main)
	}

	code := string(writer.Buffer)
	extensions := generateExtensions(writer)

	replacer := strings.NewReplacer(
		"--{extentions here}", extensions,
		"--{memory here}", fmt.Sprintf("%d", writer.Options.Memory),
		"--{code here}", code,
	)
	return []byte(replacer.Replace(luau_boilerplate))

}
