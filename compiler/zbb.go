package compiler

/** Zbb Extension - Basic Bit Manipulation */

func clz(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_COUNTLZ, src)))
}
func ctz(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_COUNTRZ, src)))
}

/* Population Count */
func cpop(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	loopBody := []*IRNode{
		// val = bit32.band(val, val - 1)
		IRStmtAssign(IRSymbol("val"),
			IRCall(BIT32_BAND, IRSymbol("val"), IRBinop("-", IRSymbol("val"), IRLit(1)))),
		IRStmtAssign(IRSymbol("count"), IRBinop("+", IRSymbol("count"), IRLit(1))),
	}
	Emit(w, IRStmtDo([]*IRNode{
		IRStmtLocal("val", "number", src),
		IRStmtLocal("count", "number", IRLit(0)),
		IRStmtWhile(IRBinop("~=", IRSymbol("val"), IRLit(0)), loopBody),
		IRStmtAssign(dst, IRSymbol("count")),
	}))
}

/* Min/Max */
func min(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(MATH_MIN, irI32(w, lhs), irI32(w, rhs))))
}
func minu(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(MATH_MIN, irU32(w, lhs), irU32(w, rhs))))
}
func max(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(MATH_MAX, irI32(w, lhs), irI32(w, rhs))))
}
func maxu(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	lhs := irArgExpr(w, command.Arguments[1])
	rhs := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(MATH_MAX, irU32(w, lhs), irU32(w, rhs))))
}

/* Sign/Zero Extension */
func sext_b(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtDo([]*IRNode{
		IRStmtLocal("val", "number", IRCall(BIT32_BAND, src, IRLitHex(0xFF))),
		IRStmtAssign(dst, IRIfExpr(
			IRBinop(">=", IRSymbol("val"), IRLitHex(0x80)),
			irI32(w, IRBinop("-", IRSymbol("val"), IRLitHex(0x100))),
			IRSymbol("val"),
		)),
	}))
}
func sext_h(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtDo([]*IRNode{
		IRStmtLocal("val", "number", IRCall(BIT32_BAND, src, IRLitHex(0xFFFF))),
		IRStmtAssign(dst, IRIfExpr(
			IRBinop(">=", IRSymbol("val"), IRLitHex(0x8000)),
			irI32(w, IRBinop("-", IRSymbol("val"), IRLitHex(0x10000))),
			IRSymbol("val"),
		)),
	}))
}
func zext_h(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BAND, src, IRLitHex(0xFFFF))))
}

/* Logical with Negate */
func andn(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BAND, a, IRCall(BIT32_BNOT, b))))
}
func orn(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BOR, a, IRCall(BIT32_BNOT, b))))
}
func xnor(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_BNOT, IRCall(BIT32_BXOR, a, b))))
}

/* Rotation */
func rol(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_LROTATE, a, IRCall(BIT32_BAND, b, IRLit(31)))))
}
func ror(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_RROTATE, a, IRCall(BIT32_BAND, b, IRLit(31)))))
}
func rori(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_RROTATE, a, b)))
}

/* Byte Operations */
func orc_b(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	// for i = 0, 3 do
	//   local byte = bit32.band(bit32.rshift(val, i*8), 0xFF)
	//   if byte ~= 0 then result = bit32.bor(result, bit32.lshift(0xFF, i*8)) end
	// end
	forBody := []*IRNode{
		IRStmtLocal("byte_", "number",
			IRCall(BIT32_BAND,
				IRCall(BIT32_RSHIFT, IRSymbol("val"), IRBinop("*", IRSymbol("i"), IRLit(8))),
				IRLitHex(0xFF))),
		IRStmtIf(
			IRBinop("~=", IRSymbol("byte_"), IRLit(0)),
			[]*IRNode{
				IRStmtAssign(IRSymbol("result"),
					IRCall(BIT32_BOR, IRSymbol("result"),
						IRCall(BIT32_LSHIFT, IRLitHex(0xFF), IRBinop("*", IRSymbol("i"), IRLit(8))))),
			},
			nil,
		),
	}
	Emit(w, IRStmtDo([]*IRNode{
		IRStmtLocal("val", "number", src),
		IRStmtLocal("result", "number", IRLit(0)),
		IRStmtForNum("i", IRLit(0), IRLit(3), forBody),
		IRStmtAssign(dst, IRSymbol("result")),
	}))
}
func rev8(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtDo([]*IRNode{
		IRStmtLocal("val", "number", src),
		IRStmtLocal("b0", "number", IRCall(BIT32_BAND, IRSymbol("val"), IRLitHex(0xFF))),
		IRStmtLocal("b1", "number", IRCall(BIT32_BAND, IRCall(BIT32_RSHIFT, IRSymbol("val"), IRLit(8)), IRLitHex(0xFF))),
		IRStmtLocal("b2", "number", IRCall(BIT32_BAND, IRCall(BIT32_RSHIFT, IRSymbol("val"), IRLit(16)), IRLitHex(0xFF))),
		IRStmtLocal("b3", "number", IRCall(BIT32_BAND, IRCall(BIT32_RSHIFT, IRSymbol("val"), IRLit(24)), IRLitHex(0xFF))),
		IRStmtAssign(dst,
			IRCall(BIT32_BOR,
				IRCall(BIT32_BOR, IRCall(BIT32_LSHIFT, IRSymbol("b0"), IRLit(24)), IRCall(BIT32_LSHIFT, IRSymbol("b1"), IRLit(16))),
				IRCall(BIT32_BOR, IRCall(BIT32_LSHIFT, IRSymbol("b2"), IRLit(8)), IRSymbol("b3")))),
	}))
}

/* Packing */
func pack(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BOR,
			IRCall(BIT32_BAND, a, IRLitHex(0xFFFF)),
			IRCall(BIT32_LSHIFT, IRCall(BIT32_BAND, b, IRLitHex(0xFFFF)), IRLit(16)))))
}
func packh(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	a   := irArgExpr(w, command.Arguments[1])
	b   := irArgExpr(w, command.Arguments[2])
	Emit(w, IRStmtAssign(dst,
		IRCall(BIT32_BOR,
			IRCall(BIT32_BAND, a, IRLitHex(0xFF)),
			IRCall(BIT32_LSHIFT, IRCall(BIT32_BAND, b, IRLitHex(0xFF)), IRLit(8)))))
}
