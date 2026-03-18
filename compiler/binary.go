package compiler

/** Binary Shifts */
func sll(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(bit32.lshift(%s, %s), 0xFFFFFFFF)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func srl(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(bit32.rshift(%s, %s), 0xFFFFFFFF)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func sra(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(bit32.arshift(%s, %s), 0xFFFFFFFF)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/* Comparision */
func slt(w *OutputWriter, command AssemblyCommand) { /* sltu & sltui instructions */
	lhs := CompileRegister(w, command.Arguments[1])
	rhs := CompileRegister(w, command.Arguments[2])
	if command.Name == "sltu" || command.Name == "sltiu" {
		WriteIndentedString(w, "%s = if (%s < %s) then 1 else 0\n", CompileRegister(w, command.Arguments[0]), wrapU32Expr(w, lhs), wrapU32Expr(w, rhs))
	} else {
		WriteIndentedString(w, "%s = if (%s < %s) then 1 else 0\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, lhs), wrapI32Expr(w, rhs))
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
	WriteIndentedString(w, "%s = bit32.band(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func xor(w *OutputWriter, command AssemblyCommand) { /* xor & xori instructions */
	WriteIndentedString(w, "%s = bit32.bxor(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func or(w *OutputWriter, command AssemblyCommand) { /* or & ori instructions */
	WriteIndentedString(w, "%s = bit32.bor(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func not(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bnot(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
