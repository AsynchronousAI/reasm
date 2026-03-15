package compiler

import log "github.com/sirupsen/logrus"

func ret(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "if r2 ~= 0 then\n")
	w.Depth++
	//WriteIndentedString(w, "print('ret', RETURN)\n")
	WriteIndentedString(w, "PC = r2\n")
	if w.Options.Trace {
		WriteIndentedString(w, "print('RET: ', PC)\n")
	}
	WriteIndentedString(w, "r2 = 0\n")
	WriteIndentedString(w, "return true\n")
	w.Depth--
	WriteIndentedString(w, "else\n")
	w.Depth++
	WriteIndentedString(w, "PC = 0\n")
	WriteIndentedString(w, "return true\n")
	w.Depth--
	WriteIndentedString(w, "end\n")
}
func call(w *OutputWriter, command AssemblyCommand) {
	var function = command.Arguments[0].Source

	/* the actual jump */
	WriteIndentedString(w, "if functions[\"%s\"] then\n", function)
	w.Depth++
	WriteIndentedString(w, "functions[\"%s\"]()\n", function)
	WriteIndentedString(w, "PC = %d\n", w.MaxPC)
	if w.Options.Trace {
		WriteIndentedString(w, "print('CALL: ', PC)\n")
	}
	WriteIndentedString(w, "return true\n")
	w.Depth--
	WriteIndentedString(w, "else\n")
	w.Depth++
	JumpTo(w, function, true)
	w.Depth--
	WriteIndentedString(w, "end\n")

	/* cut the jump */
	CutAndLink(w)
}
func move(w *OutputWriter, command AssemblyCommand) {
	WriteIndentedString(w, "%s = %s\n", CompileRegister(w, command.Arguments[0]), CompileRegister(w, command.Arguments[1]))
}

/* unimplemented */
func ebreak(w *OutputWriter, command AssemblyCommand) {
	log.Warn("EBREAK cannot be used (yet).")
}
func ecall(w *OutputWriter, command AssemblyCommand) {
	log.Warn("ECALL cannot be used (yet).")
}
func fence(w *OutputWriter, command AssemblyCommand) {
	log.Warn("FENCE cannot be used.")
}
func nop(w *OutputWriter, command AssemblyCommand) {
	if w.Options.Comments {
		WriteIndentedString(w, "-- nop\n")
	}
}
