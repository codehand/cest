package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

const newFilePerm os.FileMode = 0644

const (
	specifyFlagMessage = "Please specify either the -only, -excl, -exported, or -all flag"
	specifyFileMessage = "Please specify a file or directory containing the source"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Basepath   = filepath.Dir(b)
	Testpath   = Basepath + "/tests"
)

type OptionsCMD struct {
	OnlyFuncs          string // Regexp string for filter matches.
	ExclFuncs          string // Regexp string for excluding matches.
	ExportedFuncs      bool   // Only include exported functions.
	AllFuncs           bool   // Include all non-tested functions.
	PrintInputs        bool   // Print function parameters as part of error messages.
	Subtests           bool   // Print tests using Go 1.7 subtests
	WriteOutput        bool   // Write output to test file(s).
	TemplateDir        string // Path to custom template set
	TemplateParamsPath string // Path to custom paramters json file(s).
	OutputDir          string // Path output test
}

// Generates tests for the Go files defined in args with the given OptionsCMD.
// Logs information and errors to out. By default outputs generated tests to
// out unless specified by opt.
func Run(out io.Writer, args []string, opts *OptionsCMD) {
	if opts == nil {
		opts = &OptionsCMD{}
	}
	opt := parseOptionsCMD(out, opts)
	if opt == nil {
		return
	}
	if len(args) == 0 {
		fmt.Fprintln(out, specifyFileMessage)
		return
	}
	for _, path := range args {
		generateTests(out, path, opts.WriteOutput, opt)
	}
}

func parseOptionsCMD(out io.Writer, opt *OptionsCMD) *Options {
	if opt.OnlyFuncs == "" && opt.ExclFuncs == "" && !opt.ExportedFuncs && !opt.AllFuncs {
		fmt.Fprintln(out, specifyFlagMessage)
		return nil
	}
	onlyRE, err := parseRegexp(opt.OnlyFuncs)
	if err != nil {
		fmt.Fprintln(out, "Invalid -only regex:", err)
		return nil
	}
	exclRE, err := parseRegexp(opt.ExclFuncs)
	if err != nil {
		fmt.Fprintln(out, "Invalid -excl regex:", err)
		return nil
	}

	templateParams := map[string]interface{}{}
	jfile := opt.TemplateParamsPath
	if jfile != "" {
		buf, err := ioutil.ReadFile(jfile)
		if err != nil {
			fmt.Fprintf(out, "Failed to read from %s ,err %s", jfile, err)
			return nil
		}

		err = json.Unmarshal(buf, templateParams)
		if err != nil {
			fmt.Fprintf(out, "Failed to umarshal %s er %s", jfile, err)
			return nil
		}
	}

	if opt.OutputDir != "default" {
		// to do
	}
	return &Options{
		Only:           onlyRE,
		Exclude:        exclRE,
		Exported:       opt.ExportedFuncs,
		PrintInputs:    opt.PrintInputs,
		Subtests:       opt.Subtests,
		TemplateDir:    opt.TemplateDir,
		TemplateParams: templateParams,
		OutputDir:      opt.OutputDir,
	}
}

func parseRegexp(s string) (*regexp.Regexp, error) {
	if s == "" {
		return nil, nil
	}
	re, err := regexp.Compile(s)
	if err != nil {
		return nil, err
	}
	return re, nil
}

func generateTests(out io.Writer, path string, writeOutput bool, opt *Options) {
	// fmt.Println("generateTests:", path)
	gts, err := GenerateTests(path, opt)
	if err != nil {
		fmt.Fprintln(out, err.Error())
		return
	}
	if len(gts) == 0 {
		printAction("blue+h:black", "Skip", "No tests generated for "+path)
		return
	}
	for _, t := range gts {
		outputTest(out, t, writeOutput, opt.OutputCustomDefault())
	}
}

func outputTest(out io.Writer, t *GeneratedTest, writeOutput, defaultOutput bool) {
	if defaultOutput {
		existOrCreateDir(Testpath)
	}
	if writeOutput {
		if IsFileExist(t.Path) {
			if err := ioutil.WriteFile(t.Path, t.Output, newFilePerm); err != nil {
				fmt.Fprintln(out, err)
				return
			}
		} else {
			if err := ioutil.WriteFile(t.Path, t.Output, newFilePerm); err != nil {
				fmt.Fprintln(out, err)
				return
			}
		}
	}
	for _, t := range t.Functions {
		printAction("green+h:black", "Created", t.TestName())
	}
	if !writeOutput {
		out.Write(t.Output)
	}
}
