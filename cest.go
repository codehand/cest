package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/codehand/cest/cmd"
)

var (
	help = flag.Bool("h", false, "print this help")
	fn   = flag.String("fn", "", `regexp. generate tests for functions and methods that match func name.`)
)

var (
	onlyFuncs          = flag.String("only", "", `regexp. generate tests for functions and methods that match only. Takes precedence over -all`)
	exclFuncs          = flag.String("excl", "", `regexp. generate tests for functions and methods that don't match. Takes precedence over -only, -exported, and -all`)
	outputDir          = flag.String("output", "default", `folder output default for tests`)
	exportedFuncs      = flag.Bool("exported", false, `generate tests for exported functions and methods. Takes precedence over -only and -all`)
	allFuncs           = flag.Bool("all", false, "generate tests for all functions and methods")
	printInputs        = flag.Bool("i", false, "print test inputs in error messages")
	writeOutput        = flag.Bool("w", true, "write output to (test) files instead of stdout")
	templateDir        = flag.String("template_dir", "", `optional. Path to a directory containing custom test code templates`)
	templateParamsPath = flag.String("template_params_file", "", "read external parameters to template by json with file")
	templateParams     = flag.String("template_params", "", "read external parameters to template by json with stdin")
)

var nosubtests = true

func helper() {
	fmt.Println("USAGE")
	fmt.Println("  cest help")
	fmt.Println("")
	fmt.Println("FLAGS")
	flag.PrintDefaults()
}

func usage() string {
	return fmt.Sprintf("Usage: %s <option> <path> (try -h)", os.Args[0])
}

func main() {
	flag.Parse()
	if *help {
		helper()
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		log.Fatal(usage())
	}
	args := flag.Args()

	cmd.Run(os.Stdout, args, &cmd.OptionsCMD{
		OnlyFuncs:          *onlyFuncs,
		ExclFuncs:          *exclFuncs,
		ExportedFuncs:      *exportedFuncs,
		AllFuncs:           *allFuncs,
		PrintInputs:        *printInputs,
		Subtests:           !nosubtests,
		WriteOutput:        *writeOutput,
		TemplateDir:        *templateDir,
		TemplateParamsPath: *templateParamsPath,
		OutputDir:          *outputDir,
	})

}
