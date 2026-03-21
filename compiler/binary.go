package compiler

/** Binary Shifts */
func sll(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BAND, IRCall(BIT32_LSHIFT, lhs, rhs), IRLitHex(0xFFFFFFFF))))
}
func srl(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BAND, IRCall(BIT32_RSHIFT, lhs, rhs), IRLitHex(0xFFFFFFFF))))
}
func sra(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BAND, IRCall(BIT32_ARSHIFT, lhs, rhs), IRLitHex(0xFFFFFFFF))))
}

/* Comparison */
func slt(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	if command.Name == "sltu" || command.Name == "sltiu" {
		Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("<", irU32(w, lhs), irU32(w, rhs)), IRLit(1), IRLit(0))))
	} else {
		Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("<", irI32(w, lhs), irI32(w, rhs)), IRLit(1), IRLit(0))))
	}
}
func sgt(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	if command.Name == "sgtu" {
		Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop(">", irU32(w, lhs), irU32(w, rhs)), IRLit(1), IRLit(0))))
	} else {
		Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop(">", irI32(w, lhs), irI32(w, rhs)), IRLit(1), IRLit(0))))
	}
}
func seqz(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("==", src, IRLit(0)), IRLit(1), IRLit(0))))
}
func snez(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("~=", src, IRLit(0)), IRLit(1), IRLit(0))))
}
func sltz(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("<", src, IRLit(0)), IRLit(1), IRLit(0))))
}
func sgtz(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop(">", src, IRLit(0)), IRLit(1), IRLit(0))))
}

/** Binary Operations */
func and(w *OutputWriter, command AssemblyCommand) { /* and & andi */
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BAND, lhs, rhs)))
}
func xor(w *OutputWriter, command AssemblyCommand) { /* xor & xori */
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BXOR, lhs, rhs)))
}
func or(w *OutputWriter, command AssemblyCommand) { /* or & ori */
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BOR, lhs, rhs)))
}
func not(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BNOT, src)))
}
