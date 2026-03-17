package compiler

/** Zbs Extension - Single-Bit Instructions */

/* Bit Set */
func bset(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bor(%s, bit32.lshift(1, bit32.band(%s, 31)))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func bseti(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bor(%s, bit32.lshift(1, %s))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/* Bit Clear */
func bclr(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(%s, bit32.bnot(bit32.lshift(1, bit32.band(%s, 31))))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func bclri(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(%s, bit32.bnot(bit32.lshift(1, %s)))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/* Bit Invert */
func binv(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bxor(%s, bit32.lshift(1, bit32.band(%s, 31)))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func binvi(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bxor(%s, bit32.lshift(1, %s))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/* Bit Extract */
func bext(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(bit32.rshift(%s, bit32.band(%s, 31)), 1)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func bexti(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(bit32.rshift(%s, %s), 1)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
