package compiler

import (
	_ "embed"
	"strings"

	log "github.com/sirupsen/logrus"
)

//go:embed boilerplate.luau
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

	/* F extension */
	"fclass.s": fclass,
	"fclass.d": fclass,

	/** Arithmetic */
	"fadd.s": add,
	"fsub.s": sub,
	"fdiv.s": div,
	"fmul.s": mul,
	"fadd.d": add,
	"fsub.d": sub,
	"fdiv.d": div,
	"fmul.d": mul,
	"fneg.s": fneg,
	"fneg.d": fneg,

	/** More advanced */
	"fabs.s":  fneg,
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
	"feq.d": feq,
	"flt.d": flt,
	"fle.d": fle,

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

	"fmv.w.x":   fmv_w_x,
	"fmv.x.w":   fmv_x_w,
	"fcvt.w.s":  fcvt_w_s,
	"fcvt.wu.s": fcvt_w_s,
	"fcvt.s.w":  fcvt_s_w,
	"fcvt.s.wu": fcvt_s_w,
	"fcvt.d.s":  fcvt_d_s,
	"fcvt.w.d":  fcvt_w_d,
	"fcvt.d.w":  fcvt_d_w,

	/* Abstraction */
	"auipc": auipc,
	"ret":   ret,
	"call":  call,
	"mv":    move,
	"nop":   nop,
}
var directives = map[string]func(*OutputWriter, []string){
	".asciz":  asciz,
	".string": asciz,
	".quad":   quad,
	".word":   word,
	".byte":   byte_,
	".half":   half,
	".zero":   zero,
	".set":    set,
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
func BeforeCompilation(writer *OutputWriter) {

	/* load directives */
	WriteIndentedString(writer, "function init()\n")
	writer.Depth++
	WriteIndentedString(writer, "reset_registers()\n")

	/* pre-read the code, and write directive to top */
	for i := range writer.Commands {
		if writer.Commands[i].Type == Label {
			writer.CurrentLabel = &writer.Commands[i]
			writer.Commands[i].Ignore = true
			writer.PendingData = PendingData{} // reset so first directive under this label always saves pointer
		}
		if writer.Commands[i].Type == Instruction {
			writer.CurrentLabel.Ignore = false
		}
		if writer.Commands[i].Type == Directive {
			attributeComponents := ReadDirective(writer.Commands[i].Name)
			attributeName := attributeComponents[0]
			if _, ok := directives[attributeName]; ok {
				directives[attributeName](writer, attributeComponents)
			} else if writer.Options.Comments {
				WriteIndentedString(writer, "-- ASM DIRECTIVE: %s\n", writer.Commands[i].Name)
			}
		}
	}

	/* finish loading directives */
	WriteIndentedString(writer, "PC = %d\n", FindLabelAddress(writer, writer.Options.MainSymbol))
	WriteIndentedString(writer, "r3 = (buffer.len(memory) + %d) / 2 -- start at the center after static data\n", writer.MemoryDevelopmentPointer)
	WriteIndentedString(writer, "if r3 >= buffer.len(memory) then error(\"Not enough memory\") end\n")
	writer.Depth--
	WriteIndentedString(writer, "end\n")

	/* start code */
}
func AfterCompilation(writer *OutputWriter) []byte {
	AddEnd(writer) // end the current label, if active

	// check if invalid PC, then break
	WriteIndentedString(writer, "function start(startPosition)\n")
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
	writer.Depth--
	WriteIndentedString(writer, "end\n")

	// final code based on mode
	main := FindLabelAddress(writer, writer.Options.MainSymbol)
	if writer.Options.Mode == "main" {
		WriteString(writer, "init()\nstart(%d)\n", main)
	} else if writer.Options.Mode == "module" {
		WriteString(writer, `return setmetatable({
	init = init,
	memory = memory,
	functions = functions,
	util = {
		get_args = get_args,
		push_args = push_args,
		get_f_args = get_f_args,
		push_f_args = push_f_args,
		int_to_float = int_to_float,
		float_to_int = float_to_int,
		hi = hi,
		lo = lo,
		read_string = read_string,
		format_string = format_string,
		malloc = malloc,
		float_to_double = float_to_double,
		two_words_to_double = two_words_to_double,
		fclass = fclass,
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
		"--{code here}", code,
	)
	return []byte(replacer.Replace(luau_boilerplate))

}
