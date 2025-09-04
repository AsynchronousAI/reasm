package compiler

import (
	"strconv"
	"strings"
)

func label(w *OutputWriter, command AssemblyCommand) {
	/* end previous label */
	AddEnd(w)

	/* define it */
	w.CurrentLabel = &command
	if w.CurrentLabel.Ignore {
		return
	}

	WriteIndentedString(w, "FUNCS[%d] = function() -- %s\n", w.MaxPC, command.Name)
	w.Depth++
	w.MaxPC++
}

func save_pointer_at(w *OutputWriter, what string, where int32) {
	w.MemoryMap[what] = int(where)
}
func save_pointer(w *OutputWriter) {
	save_pointer_at(w, w.CurrentLabel.Name, w.MemoryDevelopmentPointer)
}

func asciz(w *OutputWriter, components []string) {
	var data = strings.Trim(components[1], "\"")
	w.PendingData.Data = data
	w.PendingData.Type = PendingDataTypeString

	if w.Options.Comments {
		WriteIndentedString(w, "writestring(memory, %d, \"%s\\0\") -- %s\n", w.MemoryDevelopmentPointer, data, w.CurrentLabel.Name)
	} else {
		WriteIndentedString(w, "writestring(memory, %d, \"%s\\0\")\n", w.MemoryDevelopmentPointer, data)
	}

	save_pointer(w)
	w.MemoryDevelopmentPointer += int32(len(data) + 1)
}
func quad(w *OutputWriter, components []string) {
	if w.PendingData.Type != PendingDataTypeNumeric {
		w.PendingData.Data = strconv.Itoa(int(w.MemoryDevelopmentPointer))
		save_pointer(w)
	}
	w.PendingData.Type = PendingDataTypeNumeric

	val, _ := strconv.ParseInt(components[1], 0, 0)
	WriteIndentedString(w, "writei32(memory, %d, %d)\n", w.MemoryDevelopmentPointer, val&0xFFFFFFFF)
	WriteIndentedString(w, "writei32(memory, %d, %d)\n", w.MemoryDevelopmentPointer+4, val>>32)

	w.MemoryDevelopmentPointer += 8
}
func word(w *OutputWriter, components []string) {
	if w.PendingData.Type != PendingDataTypeNumeric {
		w.PendingData.Data = strconv.Itoa(int(w.MemoryDevelopmentPointer))
		save_pointer(w)
	}
	w.PendingData.Type = PendingDataTypeNumeric

	val, _ := strconv.ParseInt(components[1], 0, 0)
	WriteIndentedString(w, "writei32(memory, %d, %d)\n", w.MemoryDevelopmentPointer, val)

	w.MemoryDevelopmentPointer += 4
}
func half(w *OutputWriter, components []string) {
	if w.PendingData.Type != PendingDataTypeNumeric {
		w.PendingData.Data = strconv.Itoa(int(w.MemoryDevelopmentPointer))
		save_pointer(w)
	}
	w.PendingData.Type = PendingDataTypeNumeric

	val, _ := strconv.ParseInt(components[1], 0, 0)
	WriteIndentedString(w, "writei16(memory, %d, %d)\n", w.MemoryDevelopmentPointer, val)

	w.MemoryDevelopmentPointer += 2
}
func byte_(w *OutputWriter, components []string) { /* byte_ to avoid overlap with the type */
	if w.PendingData.Type != PendingDataTypeNumeric {
		w.PendingData.Data = strconv.Itoa(int(w.MemoryDevelopmentPointer))
		save_pointer(w)
	}
	w.PendingData.Type = PendingDataTypeNumeric

	val, _ := strconv.ParseInt(components[1], 0, 0)
	WriteIndentedString(w, "writei16(memory, %d, %d)\n", w.MemoryDevelopmentPointer, val)

	w.MemoryDevelopmentPointer += 1
}
func zero(w *OutputWriter, components []string) {
	size, _ := strconv.ParseInt(components[1], 0, 0)
	save_pointer(w)
	WriteIndentedString(w, "fill(memory, %d, 0, %d)\n", w.MemoryDevelopmentPointer, size)

	w.MemoryDevelopmentPointer += int32(size)
}
func set(w *OutputWriter, components []string) {
	save_pointer_at(w, components[1], w.MemoryDevelopmentPointer) /* todo: support offsetted .set directives */
}
