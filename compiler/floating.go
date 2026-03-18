package compiler

/** Memory */
func fld(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = readf64(memory, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func flw(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = readf32(memory, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func fsd(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "writef64(memory, %s, %s)\n", CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[0]))
}
func fsw(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "writef32(memory, %s, %s)\n", CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[0]))
}

/** Fused */
func fmadd(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = %s * %s + %s\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]), CompileRegister(w, command.Arguments[3]))
}
func fmsub(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = %s * %s - (%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]), CompileRegister(w, command.Arguments[3]))
}
func fnmadd(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = -(%s) * %s + %s\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]), CompileRegister(w, command.Arguments[3]))
}
func fnmsub(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = -(%s) * %s - (%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]), CompileRegister(w, command.Arguments[3]))
}

/** Sign */
func fsgnj(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = math.abs(%s) * math.sign(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func fsgnjn(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = math.abs(%s) * -math.sign(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func fsgnjx(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = %s * -math.sign(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/** Other math */
func fsqrt(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = math.sqrt(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func fmin(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = math.min(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func fmax(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = math.max(%s, %s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/** Comparators */
func feq(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = if %s == %s then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func flt(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = if %s < %s then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func fle(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = if %s <= %s then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}
func fgt(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = if %s > %s then 1 else 0\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]), CompileRegister(w, command.Arguments[2]))
}

/** Conversion */
func fcvt_d_s(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func fcvt_w_s(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = float_to_int(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func fcvt_s_w(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = int_to_float(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func fcvt_d_w(w *OutputWriter, command AssemblyCommand) {
	if w.Options.Comments {

		WriteIndentedString(w, "%s = int_to_float(%s) -- Double will just be a less precise float\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))

	} else {
		WriteIndentedString(w, "%s = int_to_float(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	}
}

func fcvt_d_wu(w *OutputWriter, command AssemblyCommand) {
	if w.Options.Comments {

		WriteIndentedString(w, "%s = int_to_float(u32(%s)) -- Double will just be a less precise float\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))

	} else {
		WriteIndentedString(w, "%s = int_to_float(u32(%s))\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	}
}
func fcvt_w_d(w *OutputWriter, command AssemblyCommand) {
	if w.Options.Comments {
		WriteIndentedString(w, "%s = float_to_int(%s) -- Double will just be a less precise float\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))

	} else {
		WriteIndentedString(w, "%s = float_to_int(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
	}
}

/** Move */
func fmv_w_x(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = int_to_float(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func fmv_x_w(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = float_to_int(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}

/** Classify */
func fclass(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = fclass(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}

/** Other */
func fneg(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = -%s\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
func fabs(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = math.abs(%s)\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}
