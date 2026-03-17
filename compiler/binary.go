package compiler

/** Binary Shifts */
func sll(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = band(lshift(%s, %s), 0xFFFFFFFF)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func srl(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = band(rshift(%s, %s), 0xFFFFFFFF)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func sra(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = band(arshift(%s, %s), 0xFFFFFFFF)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/* Comparision */
func slt(w *OutputWriter, command AssemblyCommand) { /* sltu & sltui instructions */
	if command.Name == "sltu" || command.Name == "sltiu" {
		WriteIndentedString(w, "%s = if (u32(%s) < u32(%s)) then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
	} else {
		WriteIndentedString(w, "%s = if (i32(%s) < i32(%s)) then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
	}
}
func seqz(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = if (%s == 0) then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func snez(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = if (%s ~= 0) then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func sltz(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = if (%s < 0) then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func sgtz(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = if (%s > 0) then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}

/** Binary Operations */
func and(w *OutputWriter, command AssemblyCommand) { /* and & andi instructions */
	WriteIndentedString(w, "%s = band(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func xor(w *OutputWriter, command AssemblyCommand) { /* xor & xori instructions */
	WriteIndentedString(w, "%s = bxor(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func or(w *OutputWriter, command AssemblyCommand) { /* or & ori instructions */
	WriteIndentedString(w, "%s = bor(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func not(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bnot(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
