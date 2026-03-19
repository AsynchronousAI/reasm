package compiler

import "strings"

/** Memory */
func fld(w *OutputWriter, command AssemblyCommand) {
	dst  := irArgExpr(w, command.Arguments[0])
	addr := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BUFFER_READF64, IRSymbol(SYM_MEMORY), addr)))
}
func flw(w *OutputWriter, command AssemblyCommand) {
	dst  := irArgExpr(w, command.Arguments[0])
	addr := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BUFFER_READF32, IRSymbol(SYM_MEMORY), addr)))
}
func fsd(w *OutputWriter, command AssemblyCommand) {
	addr := irArgExpr(w, command.Arguments[1])
	val  := irArgExpr(w, command.Arguments[0])
	Emit(w, IRStmtCall(BUFFER_WRITEF64, IRSymbol(SYM_MEMORY), addr, val))
}
func fsw(w *OutputWriter, command AssemblyCommand) {
	addr := irArgExpr(w, command.Arguments[1])
	val  := irArgExpr(w, command.Arguments[0])
	Emit(w, IRStmtCall(BUFFER_WRITEF32, IRSymbol(SYM_MEMORY), addr, val))
}

/** Fused multiply-add family */
func fmadd(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	c   := irArgExpr(w, command.Arguments[3])
	expr := IRBinop("+", IRBinop("*", a, b), c)
	if w.Options.Accurate && strings.HasSuffix(command.Name, ".s") {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, expr)))
	} else {
		Emit(w, IRStmtAssign(dst, expr))
	}
}
func fmsub(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	c   := irArgExpr(w, command.Arguments[3])
	expr := IRBinop("-", IRBinop("*", a, b), c)
	if w.Options.Accurate && strings.HasSuffix(command.Name, ".s") {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, expr)))
	} else {
		Emit(w, IRStmtAssign(dst, expr))
	}
}
func fnmadd(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	c   := irArgExpr(w, command.Arguments[3])
	// -(a*b) + c
	expr := IRBinop("+", IRUnop("-", IRBinop("*", a, b)), c)
	if w.Options.Accurate && strings.HasSuffix(command.Name, ".s") {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, expr)))
	} else {
		Emit(w, IRStmtAssign(dst, expr))
	}
}
func fnmsub(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	c   := irArgExpr(w, command.Arguments[3])
	// -(a*b) - c
	expr := IRBinop("-", IRUnop("-", IRBinop("*", a, b)), c)
	if w.Options.Accurate && strings.HasSuffix(command.Name, ".s") {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, expr)))
	} else {
		Emit(w, IRStmtAssign(dst, expr))
	}
}

/** Sign injection */
func fsgnj(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	// math.abs(a) * math.sign(b)
	Emit(w, IRStmtAssign(dst, IRBinop("*", IRCall(MATH_ABS, a), IRCall(MATH_SIGN, b))))
}
func fsgnjn(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	// math.abs(a) * -math.sign(b)
	Emit(w, IRStmtAssign(dst, IRBinop("*", IRCall(MATH_ABS, a), IRUnop("-", IRCall(MATH_SIGN, b)))))
}
func fsgnjx(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	// a * -math.sign(b)
	Emit(w, IRStmtAssign(dst, IRBinop("*", a, IRUnop("-", IRCall(MATH_SIGN, b)))))
}

/** Other math */
func fsqrt(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(MATH_SQRT, src)))
}
func fmin(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(MATH_MIN, a, b)))
}
func fmax(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(MATH_MAX, a, b)))
}

/** Comparators */
func feq(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("==", a, b), IRLit(1), IRLit(0))))
}
func flt(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("<", a, b), IRLit(1), IRLit(0))))
}
func fle(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop("<=", a, b), IRLit(1), IRLit(0))))
}
func fgt(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop(">", a, b), IRLit(1), IRLit(0))))
}
func fge(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRIfExpr(IRBinop(">=", a, b), IRLit(1), IRLit(0))))
}

/** Float arithmetic — optional f32 rounding */
func fadd(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	expr := IRBinop("+", a, b)
	if w.Options.Accurate && strings.HasSuffix(command.Name, ".s") {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, expr)))
	} else {
		Emit(w, IRStmtAssign(dst, expr))
	}
}
func fsub(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	expr := IRBinop("-", a, b)
	if w.Options.Accurate && strings.HasSuffix(command.Name, ".s") {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, expr)))
	} else {
		Emit(w, IRStmtAssign(dst, expr))
	}
}
func fmul(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	expr := IRBinop("*", a, b)
	if w.Options.Accurate && strings.HasSuffix(command.Name, ".s") {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, expr)))
	} else {
		Emit(w, IRStmtAssign(dst, expr))
	}
}
func fdiv(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	expr := IRBinop("/", a, b)
	if w.Options.Accurate && strings.HasSuffix(command.Name, ".s") {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, expr)))
	} else {
		Emit(w, IRStmtAssign(dst, expr))
	}
}

/** Conversion */
func fcvt_d_s(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	if w.Options.Accurate {
		Emit(w, IRStmtAssign(dst, IRCast(RT_F32, src)))
	} else {
		Emit(w, IRStmtAssign(dst, src))
	}
}
func fcvt_w_s(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	// do local v = src; dst = if v>=0 then floor(v) else ceil(v); end
	Emit(w, IRStmtDo([]*IRNode{
		IRStmtLocal("v", "number", src),
		IRStmtAssign(dst, IRIfExpr(
			IRBinop(">=", IRSymbol("v"), IRLit(0)),
			IRCall(MATH_FLOOR, IRSymbol("v")),
			IRCall(MATH_CEIL, IRSymbol("v")),
		)),
	}))
}
func fcvt_s_w(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, irI32(w, src)))
}
func fcvt_s_wu(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, irU32(w, src)))
}
func fcvt_d_w(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, irI32(w, src)))
}
func fcvt_d_wu(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, irU32(w, src)))
}
func fcvt_w_d(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtDo([]*IRNode{
		IRStmtLocal("v", "number", src),
		IRStmtAssign(dst, IRIfExpr(
			IRBinop(">=", IRSymbol("v"), IRLit(0)),
			IRCall(MATH_FLOOR, IRSymbol("v")),
			IRCall(MATH_CEIL, IRSymbol("v")),
		)),
	}))
}

/** Move */
func fmv_w_x(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(RT_INT_TO_FLOAT, src)))
}
func fmv_x_w(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(RT_FLOAT_TO_INT, src)))
}

/** fflags */
func frflags(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BAND, IRSymbol(SYM_FFLAGS), IRLitHex(0x1F))))
}
func fsflags(w *OutputWriter, command AssemblyCommand) {
	if len(command.Arguments) == 1 {
		src := irArgExpr(w, command.Arguments[0])
		Emit(w, IRStmtAssign(IRSymbol(SYM_FFLAGS), IRCall(BIT32_BAND, src, IRLitHex(0x1F))))
		return
	}
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtDo([]*IRNode{
		IRStmtLocal("oldFlags", "number", IRCall(BIT32_BAND, IRSymbol(SYM_FFLAGS), IRLitHex(0x1F))),
		IRStmtAssign(IRSymbol(SYM_FFLAGS), IRCall(BIT32_BAND, src, IRLitHex(0x1F))),
		IRStmtAssign(dst, IRSymbol("oldFlags")),
	}))
}

/** Classify */
func fclass(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(RT_FCLASS, src)))
}

/** Other */
func fneg(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRUnop("-", src)))
}
func fabs(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(MATH_ABS, src)))
}
