package compiler

/** Jump */
func jump(w *OutputWriter, command AssemblyCommand) {
	JumpTo(w, command.Arguments[0].Source, false)
}
func jal(w *OutputWriter, command AssemblyCommand) {
	JumpTo(w, command.Arguments[0].Source, true)
	CutAndLink(w)
}
func jalr(w *OutputWriter, command AssemblyCommand) {
	returnReg := irArgExpr(w, command.Arguments[0])

	sourceExpr := returnReg
	if len(command.Arguments) > 1 {
		sourceExpr = irArgExpr(w, command.Arguments[1])
	}
	offsetExpr := IRLit(0)
	if len(command.Arguments) > 2 {
		offsetExpr = irArgExpr(w, command.Arguments[2])
	}

	body := []*IRNode{
		IRStmtAssign(returnReg, IRSymbol(SYM_PC)),
		IRStmtSetPC(IRBinop("+", sourceExpr, offsetExpr)),
	}
	if w.Options.Trace {
		body = append(body, IRStmtCall("print", IRLit(`"JALR: "`), IRSymbol(SYM_PC)))
	}
	body = append(body, IRStmtReturn(true))

	Emit(w, IRStmtDo(body))
	AddEnd(w)
	WriteIndentedString(w, "FUNCS[%d] = function(): boolean -- %s (extended) \n", w.MaxPC, w.CurrentLabel.Name)
	w.Depth++
	w.MaxPC++
	w.CurrentLabel.Name = IncrementFunctionName(w.CurrentLabel.Name)
}
func jr(w *OutputWriter, command AssemblyCommand) {
	src := irArgExpr(w, command.Arguments[0])
	body := []*IRNode{IRStmtSetPC(src)}
	if w.Options.Trace {
		body = append(body, IRStmtCall("print", IRLit(`"JR: "`), IRSymbol(SYM_PC)))
	}
	body = append(body, IRStmtReturn(true))
	Emit(w, IRStmtDo(body))
}

/** Branching helpers */
func irBranch(w *OutputWriter, cond *IRNode, label string) {
	JumpToIR(w, cond, label)
}

func blt(w *OutputWriter, command AssemblyCommand) {
	lhs := irArgExpr(w, command.Arguments[0])
	rhs := irArgExpr(w, command.Arguments[1])
	var cond *IRNode
	if command.Name == "bltu" {
		cond = IRBinop("<", irU32(w, lhs), irU32(w, rhs))
	} else {
		cond = IRBinop("<", irI32(w, lhs), irI32(w, rhs))
	}
	irBranch(w, cond, command.Arguments[2].Source)
}
func bnez(w *OutputWriter, command AssemblyCommand) {
	src := irArgExpr(w, command.Arguments[0])
	irBranch(w, IRBinop("~=", src, IRLit(0)), command.Arguments[1].Source)
}
func bne(w *OutputWriter, command AssemblyCommand) {
	lhs := irArgExpr(w, command.Arguments[0])
	rhs := irArgExpr(w, command.Arguments[1])
	irBranch(w, IRBinop("~=", lhs, rhs), command.Arguments[2].Source)
}
func bge(w *OutputWriter, command AssemblyCommand) {
	lhs := irArgExpr(w, command.Arguments[0])
	rhs := irArgExpr(w, command.Arguments[1])
	var cond *IRNode
	if command.Name == "bgeu" {
		cond = IRBinop(">=", irU32(w, lhs), irU32(w, rhs))
	} else {
		cond = IRBinop(">=", irI32(w, lhs), irI32(w, rhs))
	}
	irBranch(w, cond, command.Arguments[2].Source)
}
func beqz(w *OutputWriter, command AssemblyCommand) {
	src := irArgExpr(w, command.Arguments[0])
	irBranch(w, IRBinop("==", src, IRLit(0)), command.Arguments[1].Source)
}
func beq(w *OutputWriter, command AssemblyCommand) {
	lhs := irArgExpr(w, command.Arguments[0])
	rhs := irArgExpr(w, command.Arguments[1])
	irBranch(w, IRBinop("==", lhs, rhs), command.Arguments[2].Source)
}
func bgt(w *OutputWriter, command AssemblyCommand) {
	lhs := irArgExpr(w, command.Arguments[0])
	rhs := irArgExpr(w, command.Arguments[1])
	var cond *IRNode
	if command.Name == "bgtu" {
		cond = IRBinop(">", irU32(w, lhs), irU32(w, rhs))
	} else {
		cond = IRBinop(">", irI32(w, lhs), irI32(w, rhs))
	}
	irBranch(w, cond, command.Arguments[2].Source)
}
func ble(w *OutputWriter, command AssemblyCommand) {
	lhs := irArgExpr(w, command.Arguments[0])
	rhs := irArgExpr(w, command.Arguments[1])
	var cond *IRNode
	if command.Name == "bleu" {
		cond = IRBinop("<=", irU32(w, lhs), irU32(w, rhs))
	} else {
		cond = IRBinop("<=", irI32(w, lhs), irI32(w, rhs))
	}
	irBranch(w, cond, command.Arguments[2].Source)
}
func bltz(w *OutputWriter, command AssemblyCommand) {
	src := irArgExpr(w, command.Arguments[0])
	irBranch(w, IRBinop("<", src, IRLit(0)), command.Arguments[1].Source)
}
func bgtz(w *OutputWriter, command AssemblyCommand) {
	src := irArgExpr(w, command.Arguments[0])
	irBranch(w, IRBinop(">", src, IRLit(0)), command.Arguments[1].Source)
}
func blez(w *OutputWriter, command AssemblyCommand) {
	src := irArgExpr(w, command.Arguments[0])
	irBranch(w, IRBinop("<=", src, IRLit(0)), command.Arguments[1].Source)
}
func bgez(w *OutputWriter, command AssemblyCommand) {
	src := irArgExpr(w, command.Arguments[0])
	irBranch(w, IRBinop(">=", src, IRLit(0)), command.Arguments[1].Source)
}

/* AUIPC */
func auipc(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRBinop("+", IRSymbol(SYM_PC), src)))
}
