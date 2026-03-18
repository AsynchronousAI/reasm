package compiler

import (
	_ "embed"
	"strings"

	"github.com/sirupsen/logrus"
)

//go:embed math.luau
var math_extension string

var extensions = map[string]string{
	"math": math_extension,
}

func generateExtensions(writer *OutputWriter) string {
	var sb strings.Builder
	includesMath := false

	for _, name := range writer.Options.Imports {
		if name == "math" {
			includesMath = true
		}
		if ext, ok := extensions[name]; ok {
			sb.WriteString(ext)
			sb.WriteString("\n")
		} else {
			logrus.Warnf("unknown import: %s", name)
		}
	}

	if !includesMath {
		sb.WriteString(math_extension)
		sb.WriteString("\n")
	}

	return sb.String()
}
