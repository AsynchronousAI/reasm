package compiler

/** Zbb Extension - Basic Bit Manipulation */

/* Count Leading/Trailing Zeros */
func clz(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.countlz(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}

func ctz(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.countrz(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}

/* Population Count */
func cpop(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "do\n")
	w.Depth++
	WriteIndentedString(w, "local val: number = %s\n", CompileRegister(w, command.Arguments[1]))
	WriteIndentedString(w, "local count: number = 0\n")
	WriteIndentedString(w, "while val ~= 0 do\n")
	w.Depth++
	WriteIndentedString(w, "val = bit32.band(val, val - 1) -- pop LSB\n")
	WriteIndentedString(w, "count = count + 1\n")
	w.Depth--
	WriteIndentedString(w, "end\n")
	WriteIndentedString(w, "%s = count\n", CompileRegister(w, command.Arguments[0]))
	w.Depth--
	WriteIndentedString(w, "end\n")
}

/* Min/Max */
func min(w *OutputWriter, command AssemblyCommand) {
	lhs := CompileRegister(w, command.Arguments[1])
	rhs := CompileRegister(w, command.Arguments[2])
	WriteIndentedString(w, "%s = math.min(%s, %s)\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, lhs), wrapI32Expr(w, rhs))
}

func minu(w *OutputWriter, command AssemblyCommand) {
	lhs := CompileRegister(w, command.Arguments[1])
	rhs := CompileRegister(w, command.Arguments[2])
	WriteIndentedString(w, "%s = math.min(%s, %s)\n", CompileRegister(w, command.Arguments[0]), wrapU32Expr(w, lhs), wrapU32Expr(w, rhs))
}

func max(w *OutputWriter, command AssemblyCommand) {
	lhs := CompileRegister(w, command.Arguments[1])
	rhs := CompileRegister(w, command.Arguments[2])
	WriteIndentedString(w, "%s = math.max(%s, %s)\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, lhs), wrapI32Expr(w, rhs))
}

func maxu(w *OutputWriter, command AssemblyCommand) {
	lhs := CompileRegister(w, command.Arguments[1])
	rhs := CompileRegister(w, command.Arguments[2])
	WriteIndentedString(w, "%s = math.max(%s, %s)\n", CompileRegister(w, command.Arguments[0]), wrapU32Expr(w, lhs), wrapU32Expr(w, rhs))
}

/* Sign/Zero Extension */
func sext_b(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "do\n")
	w.Depth++
	WriteIndentedString(w, "local val: number = bit32.band(%s, 0xFF)\n", CompileRegister(w, command.Arguments[1]))
	WriteIndentedString(w, "%s = if val >= 0x80 then %s else val\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, "val - 0x100"))
	w.Depth--
	WriteIndentedString(w, "end\n")
}

func sext_h(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "do\n")
	w.Depth++
	WriteIndentedString(w, "local val: number = bit32.band(%s, 0xFFFF)\n", CompileRegister(w, command.Arguments[1]))
	WriteIndentedString(w, "%s = if val >= 0x8000 then %s else val\n", CompileRegister(w, command.Arguments[0]), wrapI32Expr(w, "val - 0x10000"))
	w.Depth--
	WriteIndentedString(w, "end\n")
}

func zext_h(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(%s, 0xFFFF)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}

/* Logical with Negate */
func andn(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.band(%s, bit32.bnot(%s))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func orn(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bor(%s, bit32.bnot(%s))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func xnor(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bnot(bit32.bxor(%s, %s))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/* Rotation */
func rol(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.lrotate(%s, bit32.band(%s, 31))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func ror(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.rrotate(%s, bit32.band(%s, 31))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func rori(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.rrotate(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/* Byte Operations */
func orc_b(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "do\n")
	w.Depth++
	WriteIndentedString(w, "local val: number = %s\n", CompileRegister(w, command.Arguments[1]))
	WriteIndentedString(w, "local result: number = 0\n")
	WriteIndentedString(w, "for i = 0, 3 do\n")
	w.Depth++
	WriteIndentedString(w, "local byte: number = bit32.band(bit32.rshift(val, i * 8), 0xFF)\n")
	WriteIndentedString(w, "if byte ~= 0 then\n")
	w.Depth++
	WriteIndentedString(w, "result = bit32.bor(result, bit32.lshift(0xFF, i * 8))\n")
	w.Depth--
	WriteIndentedString(w, "end\n")
	w.Depth--
	WriteIndentedString(w, "end\n")
	WriteIndentedString(w, "%s = result\n", CompileRegister(w, command.Arguments[0]))
	w.Depth--
	WriteIndentedString(w, "end\n")
}

func rev8(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "do\n")
	w.Depth++
	WriteIndentedString(w, "local val: number = %s\n", CompileRegister(w, command.Arguments[1]))
	WriteIndentedString(w, "local b0: number = bit32.band(val, 0xFF)\n")
	WriteIndentedString(w, "local b1: number = bit32.band(bit32.rshift(val, 8), 0xFF)\n")
	WriteIndentedString(w, "local b2: number = bit32.band(bit32.rshift(val, 16), 0xFF)\n")
	WriteIndentedString(w, "local b3: number = bit32.band(bit32.rshift(val, 24), 0xFF)\n")
	WriteIndentedString(w, "%s = bit32.bor(bit32.bor(bit32.lshift(b0, 24), bit32.lshift(b1, 16)), bit32.bor(bit32.lshift(b2, 8), b3))\n", CompileRegister(w, command.Arguments[0]))
	w.Depth--
	WriteIndentedString(w, "end\n")
}

/* Packing */
func pack(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bor(bit32.band(%s, 0xFFFF), bit32.lshift(bit32.band(%s, 0xFFFF), 16))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

func packh(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = bit32.bor(bit32.band(%s, 0xFF), bit32.lshift(bit32.band(%s, 0xFF), 8))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
