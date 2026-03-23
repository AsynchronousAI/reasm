package compiler

import (
	"encoding/base64"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func label(w *OutputWriter, command AssemblyCommand) {
	/* end previous label */
	AddEnd(w)

	/* define it */
	w.CurrentLabel = &command
	if w.CurrentLabel.Ignore {
		return
	}

	w.MemoryMap[command.Name] = w.MaxPC

	Emit(w, IRStmtFuncBegin(w.MaxPC, command.Name))
	w.MaxPC++
}

func asciz(w *OutputWriter, components []string) {
	if len(components) < 2 {
		return
	}

	data, err := UnescapeDirectiveString(components[1])
	if err != nil {
		log.Warnf("failed to parse string directive %q: %v", components[1], err)
		data = components[1]
	}

	w.PendingData.Data = data
	w.PendingData.Type = PendingDataTypeString

	dataWithNull := append([]byte(data), 0)
	Emit(w, IRStmtCall(BUFFER_WRITESTR,
		IRSymbol(SYM_MEMORY),
		IRLit(w.MemoryDevelopmentPointer),
		IRRawExpr(luauStringExpression(string(dataWithNull)))))

	w.MemoryDevelopmentPointer += int32(len(dataWithNull))
}

func base64data(w *OutputWriter, components []string) {
	if len(components) < 2 {
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(components[1])
	if err != nil {
		return
	}

	for i, b := range decoded {
		Emit(w, IRStmtCall(BUFFER_WRITEU8,
			IRSymbol(SYM_MEMORY),
			IRLit(int(w.MemoryDevelopmentPointer)+i),
			IRLit(int(b))))
	}

	w.MemoryDevelopmentPointer += int32(len(decoded))
	w.PendingData.Type = PendingDataTypeString
}

func local(w *OutputWriter, components []string) {
	// .local symbol
	// This is a marker for the assembler, we can ignore it for now as we treat all labels as potentially local/global.
}

func comm(w *OutputWriter, components []string) {
	// .comm symbol, length, align
	if len(components) < 3 {
		return
	}
	// symbol := components[1]
	// length, _ := strconv.Atoi(components[2])
	// align is components[3] if present

	// Handled in compilation.go
}

func quad(w *OutputWriter, components []string) {
	if w.PendingData.Type != PendingDataTypeNumeric {
		w.PendingData.Data = strconv.Itoa(int(w.MemoryDevelopmentPointer))
	}
	w.PendingData.Type = PendingDataTypeNumeric

	val, err := strconv.ParseInt(components[1], 0, 64)
	if err != nil {
		if addr, ok := resolveSymbolAsPC(w, components[1]); ok {
			val = int64(addr)
		} else {
			log.Errorf("failed to parse or resolve .quad value %q", components[1])
		}
	}
	Emit(w,
		IRStmtCall(BUFFER_WRITEI32, IRSymbol(SYM_MEMORY), IRLit(w.MemoryDevelopmentPointer), IRLit(val&0xFFFFFFFF)),
		IRStmtCall(BUFFER_WRITEI32, IRSymbol(SYM_MEMORY), IRLit(w.MemoryDevelopmentPointer+4), IRLit(val>>32)),
	)
	w.MemoryDevelopmentPointer += 8
}

func word(w *OutputWriter, components []string) {
	if w.PendingData.Type != PendingDataTypeNumeric {
		w.PendingData.Data = strconv.Itoa(int(w.MemoryDevelopmentPointer))
	}
	w.PendingData.Type = PendingDataTypeNumeric

	val, err := strconv.ParseInt(components[1], 0, 64)
	if err != nil {
		if addr, ok := resolveSymbolAsPC(w, components[1]); ok {
			val = int64(addr)
		} else {
			log.Errorf("failed to parse or resolve .word value %q", components[1])
		}
	}
	Emit(w, IRStmtCall(BUFFER_WRITEI32, IRSymbol(SYM_MEMORY), IRLit(w.MemoryDevelopmentPointer), IRLit(val)))
	w.MemoryDevelopmentPointer += 4
}

func half(w *OutputWriter, components []string) {
	if w.PendingData.Type != PendingDataTypeNumeric {
		w.PendingData.Data = strconv.Itoa(int(w.MemoryDevelopmentPointer))
	}
	w.PendingData.Type = PendingDataTypeNumeric

	val, err := strconv.ParseInt(components[1], 0, 64)
	if err != nil {
		if addr, ok := resolveSymbolAsPC(w, components[1]); ok {
			val = int64(addr)
		} else {
			log.Errorf("failed to parse or resolve .half value %q", components[1])
		}
	}
	Emit(w, IRStmtCall(BUFFER_WRITEI16, IRSymbol(SYM_MEMORY), IRLit(w.MemoryDevelopmentPointer), IRLit(val)))
	w.MemoryDevelopmentPointer += 2
}

func byte_(w *OutputWriter, components []string) { /* byte_ to avoid overlap with the type */
	if w.PendingData.Type != PendingDataTypeNumeric {
		w.PendingData.Data = strconv.Itoa(int(w.MemoryDevelopmentPointer))
	}
	w.PendingData.Type = PendingDataTypeNumeric

	val, _ := strconv.ParseInt(components[1], 0, 0)
	Emit(w, IRStmtCall(BUFFER_WRITEU8, IRSymbol(SYM_MEMORY), IRLit(w.MemoryDevelopmentPointer), IRLit(val&0xFF)))
	w.MemoryDevelopmentPointer += 1
}

func zero(w *OutputWriter, components []string) {
	size, _ := strconv.ParseInt(components[1], 0, 0)
	Emit(w, IRStmtCall(BUFFER_FILL, IRSymbol(SYM_MEMORY), IRLit(w.MemoryDevelopmentPointer), IRLit(0), IRLit(size)))
	w.MemoryDevelopmentPointer += int32(size)
}

func align(w *OutputWriter, components []string) {
	pow, _ := strconv.Atoi(components[1])
	alignSize := int32(1 << pow)
	rem := w.MemoryDevelopmentPointer % alignSize
	if rem != 0 {
		pad := alignSize - rem
		Emit(w, IRStmtCall(BUFFER_FILL, IRSymbol(SYM_MEMORY), IRLit(w.MemoryDevelopmentPointer), IRLit(0), IRLit(pad)))
		w.MemoryDevelopmentPointer += pad
	}
}

func set(w *OutputWriter, components []string) {
	// handled in compilation.go or not needed for memory layout
}
