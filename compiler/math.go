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

	// RISC-V: division by zero: rem(a, 0) = a
	if command.Name == "remu" {
		u_lhs := irU32(w, lhs)
		u_rhs := irU32(w, rhs)
		Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("==", u_rhs, IRLit(0)), u_lhs, irU32(w, IRBinop("%", u_lhs, u_rhs)))))
	} else {
		i_lhs := irI32(w, lhs)
		i_rhs := irI32(w, rhs)
		Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("==", i_rhs, IRLit(0)), i_lhs, irI32(w, IRBinop("%", i_lhs, i_rhs)))))
	}
}
func neg(w *OutputWriter, command AssemblyCommand) { /* neg & negi instructions */
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, irI32(w, IRUnop("-", irI32(w, src)))))
}

func mulh(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])

	// Use math.floor( (a * b) / 2^32 ) for mulh
	// We use u32 for all operands to ensure we stay in the realm of Luau's 53-bit mantissa correctly if possible, 
	// but for full 64-bit mul we might need care.
	// RISC-V mulhu/mulh/mulhsu have different sign handling.
	
	if command.Name == "mulhu" {
		u_lhs := irU32(w, lhs)
		u_rhs := irU32(w, rhs)
		Emit(w, IRStmtAssign(dst, irU32(w, IRCall(MATH_FLOOR, IRBinop("/", IRBinop("*", u_lhs, u_rhs), IRLit(0x100000000))))))
	} else {
		// mulh (signed)
		i_lhs := irI32(w, lhs)
		i_rhs := irI32(w, rhs)
		Emit(w, IRStmtAssign(dst, irI32(w, IRCall(MATH_FLOOR, IRBinop("/", IRBinop("*", i_lhs, i_rhs), IRLit(0x100000000))))))
	}
}
