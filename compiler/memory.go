package compiler

import "github.com/sirupsen/logrus"

/** Save Memory */
func sw(w *OutputWriter, command AssemblyCommand) {
	addr := irArgExpr(w, command.Arguments[1])
	val := irArgExpr(w, command.Arguments[0])
	Emit(w, IRStmtCall(BUFFER_WRITEI32, IRSymbol(SYM_MEMORY), addr, val))
}
func sh(w *OutputWriter, command AssemblyCommand) {
	addr := irArgExpr(w, command.Arguments[1])
	val := irArgExpr(w, command.Arguments[0])
	Emit(w, IRStmtCall(BUFFER_WRITEI16, IRSymbol(SYM_MEMORY), addr, val))
}
func sb(w *OutputWriter, command AssemblyCommand) {
	addr := irArgExpr(w, command.Arguments[1])
	val := irArgExpr(w, command.Arguments[0])
	Emit(w, IRStmtCall(BUFFER_WRITEI8, IRSymbol(SYM_MEMORY), addr, val))
}

/** Load Memory */
func li(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, src))
}
func lui(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	src := irArgExpr(w, command.Arguments[1])
	if command.Arguments[1].Modifier == "hi" {
		if addr, ok := resolveSymbolAddress(w, command.Arguments[1].Source); ok {
			hi := (addr + 0x800) >> 12                // upper 20 bits rounded
			if command.Arguments[1].Source == ".LC77" {
				logrus.Infof("DEBUG: .LC77 addr=%d hi=%d", addr, hi)
			}
			Emit(w, IRStmtAssign(dst, IRLit(hi<<12))) // shift into LUI position
			return
		} else {
			logrus.Errorf("Could not resolve as symbol or register: %s\n", command.Arguments[1].Source)
		}
	}
	Emit(w, IRStmtAssign(dst, IRCall(BIT32_LSHIFT, src, IRLit(12))))
}
func lw(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	addr := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BUFFER_READI32, IRSymbol(SYM_MEMORY), addr)))
}
func lb(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	addr := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BUFFER_READI8, IRSymbol(SYM_MEMORY), addr)))
}
func lh(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	addr := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BUFFER_READI16, IRSymbol(SYM_MEMORY), addr)))
}
func lhu(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	addr := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BUFFER_READU16, IRSymbol(SYM_MEMORY), addr)))
}
func lbu(w *OutputWriter, command AssemblyCommand) {
	dst := irArgExpr(w, command.Arguments[0])
	addr := irArgExpr(w, command.Arguments[1])
	Emit(w, IRStmtAssign(dst, IRCall(BUFFER_READU8, IRSymbol(SYM_MEMORY), addr)))
}
