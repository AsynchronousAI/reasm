package compiler

/* Math */
func add(w *OutputWriter, command AssemblyCommand) { /* add & addi instructions */
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, irI32(w, IRBinop("+", lhs, rhs))))
}
func sub(w *OutputWriter, command AssemblyCommand) { /* sub & subi instructions */
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, irI32(w, IRBinop("-", lhs, rhs))))
}
func mul(w *OutputWriter, command AssemblyCommand) { /* mul & muli instructions */
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, irI32(w, IRBinop("*", lhs, rhs))))
}
func div(w *OutputWriter, command AssemblyCommand) { /* div & divu instructions */
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	if command.Name == "divu" {
		inner := IRCall(RT_IDIV_TRUNC, irU32(w, lhs), irU32(w, rhs))
		Emit(w, IRStmtAssign(dst, irU32(w, inner)))
	} else {
		inner := IRCall(RT_IDIV_TRUNC, irI32(w, lhs), irI32(w, rhs))
		Emit(w, IRStmtAssign(dst, irI32(w, inner)))
	}
}
func rem(w *OutputWriter, command AssemblyCommand) { /* rem & remu instructions */
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	if command.Name == "remu" {
		Emit(w, IRStmtAssign(dst, irU32(w, IRBinop("%", irU32(w, lhs), irU32(w, rhs)))))
	} else {
		Emit(w, IRStmtAssign(dst, irI32(w, IRBinop("%", irI32(w, lhs), irI32(w, rhs)))))
	}
}
func neg(w *OutputWriter, command AssemblyCommand) { /* neg & negi instructions */
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, irI32(w, IRUnop("-", irI32(w, src)))))
}

/** mulh — high 32 bits of 64-bit product */
func mulh(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BAND,
			IRCall(BIT32_LSHIFT, lhs, rhs),
			IRLitHex(0xFFFFFFFF))))
}
