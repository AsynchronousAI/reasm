package compiler

import (
	_ "embed"
	"strings"

	"github.com/sirupsen/logrus"
)

//go:embed math.luau
var math_extension string

//go:embed memory.luau
var memory_extension string

//go:embed stdio.luau
var stdio_extension string

//go:embed stdlib.luau
var stdlib_extension string

var extensions = map[string]string{
	"math":   math_extension,
	"memory": memory_extension,
	"stdio":  stdio_extension,
	"stdlib": stdlib_extension,
}

func generateExtensions(writer *OutputWriter) string {
	var sb strings.Builder
	included := map[string]bool{}

	for _, name := range writer.Options.Imports {
		if ext, ok := extensions[name]; ok {
			if included[name] {
				continue
			}
			sb.WriteString(ext)
			sb.WriteString("\n")
			included[name] = true
		} else {
			logrus.Warnf("unknown import: %s", name)
		}
	}

	return sb.String()
}
