package compiler

import log "github.com/sirupsen/logrus"

func ret(w *OutputWriter, command AssemblyCommand) {
	r2 := IRSymbol(SYM_R2)
	thenBody := []*IRNode{IRStmtSetPC(r2)}
	if w.Options.Trace {
		thenBody = append(thenBody, IRStmtCall("print", IRLit(`"RET: "`), IRSymbol(SYM_PC)))
	}
	thenBody = append(thenBody,
		IRStmtAssign(r2, IRLit(0)),
		IRStmtReturn(true),
	)
	elseBody := []*IRNode{
		IRStmtSetPC(IRLit(0)),
		IRStmtReturn(true),
	}
	Emit(w, IRStmtIf(IRBinop("~=", r2, IRLit(0)), thenBody, elseBody))
}

func call(w *OutputWriter, command AssemblyCommand) {
	function := command.Arguments[0].Source
	WriteIndentedString(w, "if functions[\"%s\"] then\n", function)
	w.Depth++
	WriteIndentedString(w, "functions[\"%s\"]()\n", function)
	Emit(w, IRStmtSetPC(IRLit(w.MaxPC)))
	if w.Options.Trace {
		Emit(w, IRStmtCall("print", IRLit(`"CALL: "`), IRSymbol(SYM_PC)))
	}
	Emit(w, IRStmtReturn(true))
	w.Depth--
	WriteIndentedString(w, "else\n")
	w.Depth++
	JumpTo(w, function, true)
	w.Depth--
	WriteIndentedString(w, "end\n")
	CutAndLink(w)
}

func tail(w *OutputWriter, command AssemblyCommand) {
	function := command.Arguments[0].Source
	// tail-call: jump without saving a return address
	WriteIndentedString(w, "if functions[\"%s\"] then\n", function)
	w.Depth++
	WriteIndentedString(w, "functions[\"%s\"]()\n", function)
	r2 := IRSymbol(SYM_R2)
	thenBody := []*IRNode{
		IRStmtSetPC(r2),
		IRStmtAssign(r2, IRLit(0)),
		IRStmtReturn(true),
	}
	elseBody := []*IRNode{
		IRStmtSetPC(IRLit(0)),
		IRStmtReturn(true),
	}
	Emit(w, IRStmtIf(IRBinop("~=", r2, IRLit(0)), thenBody, elseBody))
	w.Depth--
	WriteIndentedString(w, "else\n")
	w.Depth++
	JumpTo(w, function, false)
	w.Depth--
	WriteIndentedString(w, "end\n")
	CutAndLink(w)
}

func move(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, src))
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
		Emit(w, IRStmtComment("nop"))
	}
}
