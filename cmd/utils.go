package cmd

import (
	"fmt"

	colorable "github.com/mattn/go-colorable"
	"github.com/mgutz/ansi"
)

var out = colorable.NewColorableStdout()

func printAction(color, action, target, path string) {
	if l := len(path); l > 40 {
		path = "..." + path[l-40:l]
	}
	if l := len(target); l > 25 {
		target = target[0:15] + "..."
	}
	if l := len(action); l > 8 {
		target = target[0:8] + ".."
	}

	fmt.Fprintf(out, "	%-10s  %-28s  \t%43s\n", ansi.Color(action, color), target, path)
}
