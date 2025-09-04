package compiler

/** Jump */
func jump(w *OutputWriter, command AssemblyCommand) { /* j instructions */
	JumpTo(w, command.Arguments[0].Source, false)
}
func jal(w *OutputWriter, command AssemblyCommand) { /* jal instructions */
	JumpTo(w, command.Arguments[0].Source, true)
	CutAndLink(w)
}
func jalr(w *OutputWriter, command AssemblyCommand) {
	returnAddress := CompileRegister(w, command.Arguments[0])

	sourceRegister := returnAddress
	if len(command.Arguments) > 1 {
		CompileRegister(w, command.Arguments[1])
	}
	var offset string = "0"
	if len(command.Arguments) > 2 {
		offset = CompileRegister(w, command.Arguments[2])
	}

	WriteIndentedString(w, "do\n") // wrap with a do so luau does not complain if any code is after the continue
	w.Depth++
	WriteIndentedString(w, "%s = PC\n", returnAddress)
	WriteIndentedString(w, "PC = %s + %s\n", sourceRegister, offset)

	if w.Options.Trace {
		WriteIndentedString(w, "print('JALR: ', PC)\n")
	}

	WriteIndentedString(w, "return true\n")
	w.Depth--
	WriteIndentedString(w, "end\n")
	AddEnd(w)
	WriteIndentedString(w, "FUNCS[%d] = function() -- %s (extended) \n", w.MaxPC, w.CurrentLabel.Name)
	w.Depth++
	w.MaxPC++
	w.CurrentLabel.Name = IncrementFunctionName(w.CurrentLabel.Name)
}
func jr(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "do\n") // wrap with a do so luau does not complain if any code is after the continue
	w.Depth++
	WriteIndentedString(w, "PC = %s\n", CompileRegister(w, command.Arguments[0]))
	if w.Options.Trace {
		WriteIndentedString(w, "print('JR: ', PC)\n")
	}
	WriteIndentedString(w, "return true\n")
	w.Depth--
	WriteIndentedString(w, "end\n")
}

/** Branching */
func blt(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s < %s then\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	w.Depth++
	JumpTo(w, command.Arguments[2].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func bnez(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s ~= 0 then\n", CompileRegister(w, command.Arguments[0]))
	w.Depth++
	JumpTo(w, command.Arguments[1].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func bne(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s ~= %s then\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	w.Depth++
	JumpTo(w, command.Arguments[2].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func bge(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s >= %s then\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	w.Depth++
	JumpTo(w, command.Arguments[2].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func beqz(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s == 0 then\n", CompileRegister(w, command.Arguments[0]))
	w.Depth++
	JumpTo(w, command.Arguments[1].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func beq(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s == %s then\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	w.Depth++
	JumpTo(w, command.Arguments[2].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func bgt(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s > %s then\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	w.Depth++
	JumpTo(w, command.Arguments[2].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func ble(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s <= %s then\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	w.Depth++
	JumpTo(w, command.Arguments[2].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func bltz(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s < 0 then\n", CompileRegister(w, command.Arguments[0]))
	w.Depth++
	JumpTo(w, command.Arguments[1].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func bgtz(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s > 0 then\n", CompileRegister(w, command.Arguments[0]))
	w.Depth++
	JumpTo(w, command.Arguments[1].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func blez(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s <= 0 then\n", CompileRegister(w, command.Arguments[0]))
	w.Depth++
	JumpTo(w, command.Arguments[1].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func bgez(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if %s >= 0 then\n", CompileRegister(w, command.Arguments[0]))
	w.Depth++
	JumpTo(w, command.Arguments[1].Source, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
}

/* AUIPC */
func auipc(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = PC + %s\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
