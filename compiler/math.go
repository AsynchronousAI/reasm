package compiler

/* Math */
func add(w *OutputWriter, command AssemblyCommand) { /* add & addi instructions */
	expr := CompileRegister(w, command.Arguments[1]) + " + " + CompileRegister(w, command.Arguments[2])
	WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, expr))
}
func sub(w *OutputWriter, command AssemblyCommand) { /* sub & subi instructions */
	expr := CompileRegister(w, command.Arguments[1]) + " - " + CompileRegister(w, command.Arguments[2])
	WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, expr))
}
func mul(w *OutputWriter, command AssemblyCommand) { /* mul & muli instructions */
	expr := CompileRegister(w, command.Arguments[1]) + " * " + CompileRegister(w, command.Arguments[2])
	WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, expr))
}
func div(w *OutputWriter, command AssemblyCommand) { /* div & divi instructions */
	lhs := CompileRegister(w, command.Arguments[1])
	rhs := CompileRegister(w, command.Arguments[2])
	if command.Name == "divu" {
		expr := "idiv_trunc(" + wrapU32Expr(w, lhs) + ", " + wrapU32Expr(w, rhs) + ")"
		WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), wrapU32Expr(w, expr))
	} else {
		expr := "idiv_trunc(" + wrapI32Expr(w, lhs) + ", " + wrapI32Expr(w, rhs) + ")"
		WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, expr))
	}
}
func rem(w *OutputWriter, command AssemblyCommand) { /* rem & remi instructions */
	lhs := CompileRegister(w, command.Arguments[1])
	rhs := CompileRegister(w, command.Arguments[2])
	if command.Name == "remu" {
		expr := wrapU32Expr(w, lhs) + " % " + wrapU32Expr(w, rhs)
		WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), wrapU32Expr(w, expr))
	} else {
		expr := wrapI32Expr(w, lhs) + " % " + wrapI32Expr(w, rhs)
		WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, expr))
	}
}
func neg(w *OutputWriter, command AssemblyCommand) { /* neg & negi instructions */
	expr := "-" + wrapI32Expr(w, CompileRegister(w, command.Arguments[1]))
	WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, expr))
}

/** Math Descendants */
func mulh(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(bit32.lshift(%s, %s), 0xFFFFFFFF)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
