package compiler

import (
	_ "embed"
	"strings"

	"github.com/sirupsen/logrus"
)

//go:embed templates/math.luau
var math_extension string

//go:embed templates/memory.luau
var memory_extension string

//go:embed templates/string.luau
var string_extension string

//go:embed templates/ctype.luau
var ctype_extension string

//go:embed templates/stdio.luau
var stdio_extension string

//go:embed templates/stdlib.luau
var stdlib_extension string

//go:embed templates/time.luau
var time_extension string

var extensions = map[string]string{
	"math":   math_extension,
	"memory": memory_extension,
	"string": string_extension,
	"ctype":  ctype_extension,
	"stdio":  stdio_extension,
	"stdlib": stdlib_extension,
	"time":   time_extension,
}

func generateExtensions(writer *OutputWriter) string {
	var sb strings.Builder
	included := map[string]bool{}

	for _, name := range writer.Options.Imports {
		if ext, ok := extensions[name]; ok || name == "all" {
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

	defaultLibs := []string{"memory", "string", "ctype", "stdlib", "stdio", "time"}
	for _, name := range defaultLibs {
		if ext, ok := extensions[name]; ok && included[name] || included["all"] {
			sb.WriteString(ext)
			sb.WriteString("\n")
			included[name] = true
		}
	}

	return sb.String()
}
