package compiler

/** Save Memory */
func sw(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "writei32(memory, %s, %s)\n", CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[0]))
}
func sh(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "writei16(memory, %s, %s)\n", CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[0]))
}
func sb(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "writei8(memory, %s, %s)\n", CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[0]))
}

/** Load Memory */
func li(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func lui(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.lshift(%s, 12)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func lw(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = readi32(memory, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func lb(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = readi8(memory, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func lh(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = readi16(memory, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func lhu(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = readu16(memory, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func lbu(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = readu8(memory, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
