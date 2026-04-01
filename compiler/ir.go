package compiler

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ---------------------------------------------------------------------------
// IR Node kinds
// ---------------------------------------------------------------------------
type IRKind uint8

const (
	IRAssign    IRKind = iota // dst = expr
	IRIf                      // if cond { body... }
	IRDo                      // do { body... }
	IRReturn                  // return <bool>
	IRSetPC                   // PC = expr
	IRSetReg                  // rN = expr   (alias for IRAssign targeting a register)
	IRComment                 // -- text
	IRRaw                     // raw Luau string (escape hatch for things not worth modelling)
	IRFuncBegin               // FUNCS[n] = function(): boolean -- label
	IRFuncEnd                 // end  (closes a FUNCS entry)
	IRLocalDecl               // local name: type = expr
	IRWhile                   // while cond { body... }
	IRForNum                  // for i = start, limit { body... }
	IRPCInc                   // PC += 1
	IRError                   // error("…")
	IRFuncCall                // standalone function call statement: fn(args…)

	IRExprReg    // rN
	IRExprLit    // numeric / string literal
	IRExprSym    // symbol name (label, global)
	IRExprBinop  // lhs op rhs
	IRExprUnop   // op operand
	IRExprCall   // fn(args…)
	IRExprIfExpr // if cond then a else b
	IRExprCast   // i32(x) / u32(x) / f32(x)
	IRExprIndex  // table[key]
	IRExprField  // table.field
)

var irKindNames = [...]string{
	"IRAssign",
	"IRIf",
	"IRDo",
	"IRReturn",
	"IRSetPC",
	"IRSetReg",
	"IRComment",
	"IRRaw",
	"IRFuncBegin",
	"IRFuncEnd",
	"IRLocalDecl",
	"IRWhile",
	"IRForNum",
	"IRPCInc",
	"IRError",
	"IRFuncCall",
	"IRExprReg",
	"IRExprLit",
	"IRExprSym",
	"IRExprBinop",
	"IRExprUnop",
	"IRExprCall",
	"IRExprIfExpr",
	"IRExprCast",
	"IRExprIndex",
	"IRExprField",
}

const (
	BIT32_BAND    = "bit32.band"
	BIT32_BOR     = "bit32.bor"
	BIT32_BXOR    = "bit32.bxor"
	BIT32_BNOT    = "bit32.bnot"
	BIT32_LSHIFT  = "bit32.lshift"
	BIT32_RSHIFT  = "bit32.rshift"
	BIT32_ARSHIFT = "bit32.arshift"
	BIT32_LROTATE = "bit32.lrotate"
	BIT32_RROTATE = "bit32.rrotate"
	BIT32_COUNTLZ = "bit32.countlz"
	BIT32_COUNTRZ = "bit32.countrz"

	MATH_ABS   = "math.abs"
	MATH_SQRT  = "math.sqrt"
	MATH_MIN   = "math.min"
	MATH_MAX   = "math.max"
	MATH_SIGN  = "math.sign"
	MATH_FLOOR = "math.floor"
	MATH_CEIL  = "math.ceil"

	BUFFER_READI8   = "buffer.readi8"
	BUFFER_READI16  = "buffer.readi16"
	BUFFER_READI32  = "buffer.readi32"
	BUFFER_READU8   = "buffer.readu8"
	BUFFER_READU16  = "buffer.readu16"
	BUFFER_READU32  = "buffer.readu32"
	BUFFER_READF32  = "buffer.readf32"
	BUFFER_READF64  = "buffer.readf64"
	BUFFER_WRITEI8  = "buffer.writei8"
	BUFFER_WRITEI16 = "buffer.writei16"
	BUFFER_WRITEI32 = "buffer.writei32"
	BUFFER_WRITEU8  = "buffer.writeu8"
	BUFFER_WRITEF32 = "buffer.writef32"
	BUFFER_WRITEF64 = "buffer.writef64"
	BUFFER_WRITESTR = "buffer.writestring"
	BUFFER_FILL     = "buffer.fill"
	BUFFER_LEN      = "buffer.len"

	RT_I32          = "i32"
	RT_U32          = "u32"
	RT_F32          = "f32"
	RT_INT_TO_FLOAT = "int_to_float"
	RT_FLOAT_TO_INT = "float_to_int"
	RT_IDIV_TRUNC   = "idiv_trunc"
	RT_FCLASS       = "fclass"
	RT_RESET_REGS   = "reset_registers"
	RT_FLUSH_STDOUT = "flush_stdout"

	SYM_MEMORY    = "memory"
	SYM_PC        = "PC"
	SYM_R2        = "r2"
	SYM_FFLAGS    = "fflags"
	SYM_FUNCS     = "FUNCS"
	SYM_FUNCTIONS = "functions"
)

// ---------------------------------------------------------------------------
// IRNode
// ---------------------------------------------------------------------------
type IRNode struct {
	Kind     IRKind
	Op       string    // operator text, function name, literal value, symbol name, cast name
	Operands []*IRNode // sub-expressions
	Body     []*IRNode // then-branch / loop body / function body
	Else     []*IRNode // else branch (IRIf only)
	Comment  string    // inline comment
	IntVal   int       // used by IRFuncBegin, IRReturn (0=false,1=true), IRPCInc
	BoolVal  bool      // used by IRReturn
}

type exprResult struct {
	expr  string
	stmts []string
}

type inlineHelper struct {
	arity  int
	expand func(*OutputWriter, []string) (string, []string)
}

var inlineHelpers = map[string]inlineHelper{
	RT_U32: {1, func(_ *OutputWriter, args []string) (string, []string) {
		return fmt.Sprintf("((math.floor(%s %% 0x100000000) + 0x100000000) %% 0x100000000)", args[0]), nil
	}},
	RT_I32: {1, func(w *OutputWriter, args []string) (string, []string) {
		u32 := fmt.Sprintf("((math.floor(%s %% 0x100000000) + 0x100000000) %% 0x100000000)", args[0])
		return fmt.Sprintf("(if %s >= 0x80000000 then %s - 0x100000000 else %s)", u32, u32, u32), nil
	}},
	RT_F32: {1, func(_ *OutputWriter, args []string) (string, []string) {
		return "buffer.readf32(temp, 0)", []string{
			fmt.Sprintf("buffer.writef32(temp, 0, %s)", args[0]),
		}
	}},
	RT_INT_TO_FLOAT: {1, func(_ *OutputWriter, args []string) (string, []string) {
		return "buffer.readf32(temp, 0)", []string{
			fmt.Sprintf("buffer.writei32(temp, 0, %s)", args[0]),
		}
	}},
	RT_FLOAT_TO_INT: {1, func(_ *OutputWriter, args []string) (string, []string) {
		return "buffer.readi32(temp, 0)", []string{
			fmt.Sprintf("buffer.writef32(temp, 0, %s)", args[0]),
		}
	}},
	RT_IDIV_TRUNC: {2, func(w *OutputWriter, args []string) (string, []string) {
		a, b := args[0], args[1]
		return fmt.Sprintf("(if %s == 0 then -1 else if %s == -1 and %s == -2147483648 then -2147483648 else if %s >= 0 then (%s - (%s %% %s)) // %s else -(- %s - (- %s %% %s)) // %s)",
			b, b, a, a, a, a, b, b, a, a, b, b), nil
	}},
	RT_FCLASS: {1, func(w *OutputWriter, args []string) (string, []string) {
		x := args[0]
		return fmt.Sprintf("(if %s ~= %s then 0x200 elseif %s == math.huge then 0x80 elseif %s == -math.huge then 0x01 elseif %s == 0 then (if 1/%s == -math.huge then 0x08 else 0x10) else (if math.abs(%s) < 2.2250738585072014e-308 then (if %s > 0 then 0x20 else 0x04) else (if %s > 0 then 0x40 else 0x02)))",
			x, x, x, x, x, x, x, x, x), nil
	}},
}

func expandInlineHelper(w *OutputWriter, op string, args []string) (exprResult, bool) {
	if helper, ok := inlineHelpers[op]; ok && len(args) == helper.arity {
		expr, stmts := helper.expand(w, args)
		return exprResult{expr: expr, stmts: stmts}, true
	}
	return exprResult{}, false
}

func (k IRKind) String() string {
	idx := int(k)
	if idx >= 0 && idx < len(irKindNames) {
		return irKindNames[idx]
	}
	return fmt.Sprintf("IRKind(%d)", k)
}

func (k IRKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

// ---------------------------------------------------------------------------
// Expression constructors
// ---------------------------------------------------------------------------
func IRReg(index int) *IRNode {
	return &IRNode{Kind: IRExprReg, Op: regVarName(index)}
}
func IRRegName(name string) *IRNode {
	return &IRNode{Kind: IRExprReg, Op: name}
}
func IRLit(v interface{}) *IRNode {
	return &IRNode{Kind: IRExprLit, Op: fmt.Sprintf("%v", v)}
}
func IRLitHex(v int) *IRNode {
	return &IRNode{Kind: IRExprLit, Op: fmt.Sprintf("0x%X", v)}
}
func IRSymbol(name string) *IRNode {
	return &IRNode{Kind: IRExprSym, Op: name}
}
func IRCall(fn string, args ...*IRNode) *IRNode {
	return &IRNode{Kind: IRExprCall, Op: fn, Operands: args}
}
func IRBinop(op string, lhs, rhs *IRNode) *IRNode {
	return &IRNode{Kind: IRExprBinop, Op: op, Operands: []*IRNode{lhs, rhs}}
}
func IRUnop(op string, operand *IRNode) *IRNode {
	return &IRNode{Kind: IRExprUnop, Op: op, Operands: []*IRNode{operand}}
}
func IRIfExpr(cond, then, els *IRNode) *IRNode {
	return &IRNode{Kind: IRExprIfExpr, Operands: []*IRNode{cond, then, els}}
}
func IRCast(fn string, expr *IRNode) *IRNode { // runtime cast such as i32(x)
	return &IRNode{Kind: IRExprCast, Op: fn, Operands: []*IRNode{expr}}
}
func IRIndex(table, key *IRNode) *IRNode { // table[key]
	return &IRNode{Kind: IRExprIndex, Operands: []*IRNode{table, key}}
}
func IRField(table *IRNode, field string) *IRNode { // table.field
	return &IRNode{Kind: IRExprField, Op: field, Operands: []*IRNode{table}}
}
func IRRawExpr(text string) *IRNode {
	return &IRNode{Kind: IRExprLit, Op: text}
}

// ---------------------------------------------------------------------------
// Conditional cast helpers (respects w.Options.Accurate)
// ---------------------------------------------------------------------------
func irI32(w *OutputWriter, expr *IRNode) *IRNode {
	if !w.Options.Accurate {
		return expr
	}
	return IRCast(RT_I32, expr)
}
func irU32(w *OutputWriter, expr *IRNode) *IRNode {
	if !w.Options.Accurate {
		return expr
	}
	return IRCast(RT_U32, expr)
}

// ---------------------------------------------------------------------------
// Statement constructors
// ---------------------------------------------------------------------------
func IRStmtAssign(dst, src *IRNode) *IRNode {
	return &IRNode{Kind: IRAssign, Operands: []*IRNode{dst, src}}
}
func IRStmtAssignComment(dst, src *IRNode, comment string) *IRNode {
	n := IRStmtAssign(dst, src)
	n.Comment = comment
	return n
}
func IRStmtIf(cond *IRNode, body []*IRNode, elseBody []*IRNode) *IRNode {
	return &IRNode{Kind: IRIf, Operands: []*IRNode{cond}, Body: body, Else: elseBody}
}
func IRStmtDo(body []*IRNode) *IRNode {
	return &IRNode{Kind: IRDo, Body: body}
}
func IRStmtReturn(val bool) *IRNode {
	v := 0
	if val {
		v = 1
	}
	return &IRNode{Kind: IRReturn, BoolVal: val, IntVal: v}
}
func IRStmtSetPC(expr *IRNode) *IRNode {
	return &IRNode{Kind: IRSetPC, Operands: []*IRNode{expr}}
}
func IRStmtComment(text string) *IRNode {
	return &IRNode{Kind: IRComment, Op: text}
}
func IRStmtRaw(text string) *IRNode {
	return &IRNode{Kind: IRRaw, Op: text}
}
func IRStmtFuncBegin(pc int, label string) *IRNode {
	return &IRNode{Kind: IRFuncBegin, IntVal: pc, Op: label}
}
func IRStmtFuncEnd() *IRNode {
	return &IRNode{Kind: IRFuncEnd}
}
func IRStmtLocal(name, typ string, expr *IRNode) *IRNode {
	return &IRNode{Kind: IRLocalDecl, Op: name, Comment: typ, Operands: []*IRNode{expr}}
}
func IRStmtWhile(cond *IRNode, body []*IRNode) *IRNode {
	return &IRNode{Kind: IRWhile, Operands: []*IRNode{cond}, Body: body}
}
func IRStmtForNum(varName string, start, limit *IRNode, body []*IRNode) *IRNode {
	return &IRNode{Kind: IRForNum, Op: varName, Operands: []*IRNode{start, limit}, Body: body}
}
func IRStmtPCInc() *IRNode {
	return &IRNode{Kind: IRPCInc}
}
func IRStmtError(msg string) *IRNode {
	return &IRNode{Kind: IRError, Op: msg}
}

func IRStmtCall(fn string, args ...*IRNode) *IRNode {
	return &IRNode{Kind: IRFuncCall, Op: fn, Operands: args}
}

func appendStmts(dst []string, src []string) []string {
	if len(src) == 0 {
		return dst
	}
	return append(dst, src...)
}

func flushExprStmts(w *OutputWriter, stmts []string) {
	for _, stmt := range stmts {
		WriteIndentedString(w, "%s\n", stmt)
	}
}

func emitExpr(w *OutputWriter, n *IRNode) exprResult {
	switch n.Kind {
	case IRExprReg, IRExprSym, IRExprLit:
		return exprResult{expr: n.Op}
	case IRExprCall:
		stmts := make([]string, 0)
		args := make([]string, len(n.Operands))
		for i, operand := range n.Operands {
			sub := emitExpr(w, operand)
			stmts = appendStmts(stmts, sub.stmts)
			args[i] = sub.expr
		}
		if inline, ok := expandInlineHelper(w, n.Op, args); ok {
			inline.stmts = appendStmts(stmts, inline.stmts)
			return inline
		}
		return exprResult{
			expr:  fmt.Sprintf("%s(%s)", n.Op, strings.Join(args, ", ")),
			stmts: stmts,
		}
	case IRExprBinop:
		lhs := emitExpr(w, n.Operands[0])
		rhs := emitExpr(w, n.Operands[1])
		return exprResult{
			expr:  fmt.Sprintf("%s %s %s", lhs.expr, n.Op, rhs.expr),
			stmts: appendStmts(lhs.stmts, rhs.stmts),
		}
	case IRExprUnop:
		operand := emitExpr(w, n.Operands[0])
		return exprResult{
			expr:  fmt.Sprintf("%s(%s)", n.Op, operand.expr),
			stmts: operand.stmts,
		}
	case IRExprIfExpr:
		cond := emitExpr(w, n.Operands[0])
		thenExpr := emitExpr(w, n.Operands[1])
		elseExpr := emitExpr(w, n.Operands[2])
		stmts := appendStmts(cond.stmts, appendStmts(thenExpr.stmts, elseExpr.stmts))
		return exprResult{
			expr:  fmt.Sprintf("if %s then %s else %s", cond.expr, thenExpr.expr, elseExpr.expr),
			stmts: stmts,
		}
	case IRExprCast:
		operand := emitExpr(w, n.Operands[0])
		if inline, ok := expandInlineHelper(w, n.Op, []string{operand.expr}); ok {
			inline.stmts = appendStmts(operand.stmts, inline.stmts)
			return inline
		}
		return exprResult{
			expr:  fmt.Sprintf("%s(%s)", n.Op, operand.expr),
			stmts: operand.stmts,
		}
	case IRExprIndex:
		lhs := emitExpr(w, n.Operands[0])
		rhs := emitExpr(w, n.Operands[1])
		return exprResult{
			expr:  fmt.Sprintf("%s[%s]", lhs.expr, rhs.expr),
			stmts: appendStmts(lhs.stmts, rhs.stmts),
		}
	case IRExprField:
		base := emitExpr(w, n.Operands[0])
		return exprResult{
			expr:  fmt.Sprintf("%s.%s", base.expr, n.Op),
			stmts: base.stmts,
		}
	default:
		return exprResult{expr: "<expr?>"}
	}
}

func emitIRStatements(w *OutputWriter, nodes []*IRNode) {
	for _, n := range nodes {
		emitStmt(w, n)
	}
}
func emitStmt(w *OutputWriter, n *IRNode) {
	switch n.Kind {
	case IRAssign:
		dst := emitExpr(w, n.Operands[0])
		src := emitExpr(w, n.Operands[1])
		flushExprStmts(w, dst.stmts)
		flushExprStmts(w, src.stmts)
		if n.Comment != "" {
			WriteIndentedString(w, "%s = %s -- %s\n", dst.expr, src.expr, n.Comment)
		} else {
			WriteIndentedString(w, "%s = %s\n", dst.expr, src.expr)
		}

	case IRSetPC:
		expr := emitExpr(w, n.Operands[0])
		flushExprStmts(w, expr.stmts)
		if n.Comment != "" {
			WriteIndentedString(w, "PC = %s -- %s\n", expr.expr, n.Comment)
		} else {
			WriteIndentedString(w, "PC = %s\n", expr.expr)
		}

	case IRReturn:
		if n.BoolVal {
			WriteIndentedString(w, "return true\n")
		} else {
			WriteIndentedString(w, "return false\n")
		}

	case IRComment:
		WriteIndentedString(w, "-- %s\n", n.Op)

	case IRRaw:
		WriteIndentedString(w, "%s\n", n.Op)

	case IRError:
		WriteIndentedString(w, "error(%q)\n", n.Op)

	case IRFuncCall:
		args := make([]exprResult, len(n.Operands))
		for i, operand := range n.Operands {
			args[i] = emitExpr(w, operand)
			flushExprStmts(w, args[i].stmts)
		}
		argExprs := make([]string, len(args))
		for i, arg := range args {
			argExprs[i] = arg.expr
		}
		WriteIndentedString(w, "%s(%s)\n", n.Op, strings.Join(argExprs, ", "))

	case IRLocalDecl:
		expr := emitExpr(w, n.Operands[0])
		typ := n.Comment // reused field
		flushExprStmts(w, expr.stmts)
		if typ != "" {
			WriteIndentedString(w, "local %s: %s = %s\n", n.Op, typ, expr.expr)
		} else {
			WriteIndentedString(w, "local %s = %s\n", n.Op, expr.expr)
		}

	case IRPCInc:
		WriteIndentedString(w, "PC += 1\n")

	case IRFuncBegin:
		WriteIndentedString(w, "FUNCS[%d] = function(): boolean -- %s\n", n.IntVal, n.Op)
		w.Depth++

	case IRFuncEnd:
		w.Depth--
		if w.Options.Comments && w.CurrentLabel != nil {
			WriteIndentedString(w, "end -- %s\n", w.CurrentLabel.Name)
		} else {
			WriteIndentedString(w, "end\n")
		}

	case IRIf:
		cond := emitExpr(w, n.Operands[0])
		flushExprStmts(w, cond.stmts)
		WriteIndentedString(w, "if %s then\n", cond.expr)
		w.Depth++
		emitIRStatements(w, n.Body)
		w.Depth--
		if len(n.Else) > 0 {
			WriteIndentedString(w, "else\n")
			w.Depth++
			emitIRStatements(w, n.Else)
			w.Depth--
		}
		WriteIndentedString(w, "end\n")

	case IRDo:
		WriteIndentedString(w, "do\n")
		w.Depth++
		emitIRStatements(w, n.Body)
		w.Depth--
		WriteIndentedString(w, "end\n")

	case IRWhile:
		cond := emitExpr(w, n.Operands[0])
		flushExprStmts(w, cond.stmts)
		WriteIndentedString(w, "while %s do\n", cond.expr)
		w.Depth++
		emitIRStatements(w, n.Body)
		w.Depth--
		WriteIndentedString(w, "end\n")

	case IRForNum:
		start := emitExpr(w, n.Operands[0])
		limit := emitExpr(w, n.Operands[1])
		flushExprStmts(w, start.stmts)
		flushExprStmts(w, limit.stmts)
		WriteIndentedString(w, "for %s = %s, %s do\n", n.Op, start.expr, limit.expr)
		w.Depth++
		emitIRStatements(w, n.Body)
		w.Depth--
		WriteIndentedString(w, "end\n")

	default:
		WriteIndentedString(w, "-- [unknown IR kind %d]\n", n.Kind)
	}
}
func Emit(w *OutputWriter, nodes ...*IRNode) {
	if w.Options.LogIR {
		w.IRNodes = append(w.IRNodes, nodes...)
	}
	emitIRStatements(w, nodes)
}
func dumpIRAsJSON(w *OutputWriter) {
	if len(w.IRNodes) == 0 {
		log.Debug("IR dump requested but no nodes recorded")
		return
	}
	data, err := json.MarshalIndent(w.IRNodes, "", "  ")
	if err != nil {
		log.WithError(err).Error("failed to marshal IR nodes for logging")
		return
	}
	log.Debugf("IR JSON dump:\n%s", data)
}
