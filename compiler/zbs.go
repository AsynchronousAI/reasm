package compiler

/** Zbs Extension - Single-Bit Instructions */

/* Bit Set */
func bset(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BOR, a, IRCall(BIT32_LSHIFT, IRLit(1), IRCall(BIT32_BAND, b, IRLit(31))))))
}
func bseti(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BOR, a, IRCall(BIT32_LSHIFT, IRLit(1), b))))
}

/* Bit Clear */
func bclr(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BAND, a, IRCall(BIT32_BNOT, IRCall(BIT32_LSHIFT, IRLit(1), IRCall(BIT32_BAND, b, IRLit(31)))))))
}
func bclri(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BAND, a, IRCall(BIT32_BNOT, IRCall(BIT32_LSHIFT, IRLit(1), b)))))
}

/* Bit Invert */
func binv(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BXOR, a, IRCall(BIT32_LSHIFT, IRLit(1), IRCall(BIT32_BAND, b, IRLit(31))))))
}
func binvi(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BXOR, a, IRCall(BIT32_LSHIFT, IRLit(1), b))))
}

/* Bit Extract */
func bext(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BAND, IRCall(BIT32_RSHIFT, a, IRCall(BIT32_BAND, b, IRLit(31))), IRLit(1))))
}
func bexti(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BAND, IRCall(BIT32_RSHIFT, a, b), IRLit(1))))
}
