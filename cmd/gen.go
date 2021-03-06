package cmd

import (
	"fmt"
	"go/importer"
	"go/types"
	"path"
	"regexp"
	"sort"
	"sync"
)

type Options struct {
	Only           *regexp.Regexp
	Exclude        *regexp.Regexp
	Exported       bool
	PrintInputs    bool
	Subtests       bool
	Importer       func() types.Importer
	TemplateDir    string
	TemplateParams map[string]interface{}
	OutputDir      string
	CurrentDir     string
	IsSamePkg      bool
}

func (p *Options) OutputCustomDefault() bool {
	return p.OutputDir != "."
}

type GeneratedTest struct {
	Path      string      // The test file's absolute path.
	Functions []*Function // The functions with new test methods.
	Output    []byte      // The contents of the test file.
}

func GenerateTests(srcPath string, opt *Options) ([]*GeneratedTest, error) {
	if opt == nil {
		opt = &Options{}
	}
	if opt.Importer == nil || opt.Importer() == nil {
		opt.Importer = importer.Default
	}
	// fmt.Println("src output: ", srcPath)

	// linux support char `...`
	// unix support char `../..`

	if srcPath == "../.." || srcPath == "..." {
		srcPath = "."
		rootFiles, _, err := RootFiles(srcPath)
		if err != nil {
			panic(err)
		}
		if opt.OutputDir == "." {
			opt.IsSamePkg = true
		}
		// current is [0]

		lst := make([]*GeneratedTest, 0)
		for _, rt := range rootFiles {
			scFiles, err := Files(string(rt))
			if err != nil {
				panic(err)
			}

			scfiles, err := Files(path.Dir(string(rt)))
			if err != nil {
				panic(err)
			}

			if sfs, err := parallelize(scFiles, scfiles, opt, string(rt), opt.CurrentDir); err == nil {
				for _, item := range sfs {
					lst = append(lst, item)
				}
			}
		}
		return lst, nil
	}
	srcFiles, err := Files(srcPath)
	if err != nil {
		return nil, fmt.Errorf("Files: %v", err)
	}

	files, err := Files(path.Dir(srcPath))
	if err != nil {
		return nil, fmt.Errorf("Files: %v", err)
	}

	return parallelize(srcFiles, files, opt, srcPath, opt.CurrentDir)
}

// result stores a generateTest result.
type result struct {
	gt  *GeneratedTest
	err error
}

// parallelize generates tests for the given source files concurrently.
func parallelize(srcFiles, files []Path, opt *Options, srcPath, rootPath string) ([]*GeneratedTest, error) {
	var wg sync.WaitGroup
	rs := make(chan *result, len(srcFiles))
	for _, src := range srcFiles {
		wg.Add(1)
		// Worker
		go func(src Path) {
			defer wg.Done()
			r := &result{}
			r.gt, r.err = generateTest(src, files, opt, srcPath, rootPath)
			rs <- r
		}(src)
	}
	// Closer.
	go func() {
		wg.Wait()
		close(rs)
	}()
	return readResults(rs)
}

// readResults reads the result channel.
func readResults(rs <-chan *result) ([]*GeneratedTest, error) {
	var gts []*GeneratedTest
	for r := range rs {
		if r.err != nil {
			return nil, r.err
		}
		if r.gt != nil {
			gts = append(gts, r.gt)
		}
	}
	return gts, nil
}

func generateTest(src Path, files []Path, opt *Options, srcPath, rootPath string) (*GeneratedTest, error) {
	p := &Parser{Importer: opt.Importer()}
	sr, err := p.Parse(string(src), files)

	if err != nil {
		return nil, fmt.Errorf("Parser.Parse source file: %v", err)
	}

	h := sr.Header
	h.PkgName = h.Package
	h.Code = nil // Code is only needed from parsed test files.
	testPath := Path(src).TestPath()
	// output
	if opt.OutputCustomDefault() {
		testPath = Path(src).TestPathDefault(rootPath, opt.OutputDir)
	}

	// fmt.Println("test abc:", testPath)
	h, tf, err := parseTestFile(p, testPath, h)
	if err != nil {
		return nil, err
	}

	funcs, h := testableFuncs(h, sr.Funcs, opt.Only, opt.Exclude, opt.Exported, opt.OutputCustomDefault(), opt.IsSamePkg, tf)
	if len(funcs) == 0 || h.Package == "main" {
		return nil, nil
	}

	if opt.OutputCustomDefault() {
		h.Package = "tests" // CHECK
	}

	b, err := ProcessOutput(h, funcs, &OptionsOutput{
		PrintInputs:    opt.PrintInputs,
		Subtests:       opt.Subtests,
		TemplateDir:    opt.TemplateDir,
		TemplateParams: opt.TemplateParams,
	})
	if err != nil {
		return nil, fmt.Errorf("output.Process: %v", err)
	}

	return &GeneratedTest{
		Path:      testPath,
		Functions: funcs,
		Output:    b,
	}, nil
}

func parseTestFile(p *Parser, testPath string, h *Header) (*Header, []string, error) {
	if !IsFileExist(testPath) {
		// fmt.Println("return if not exist")
		return h, nil, nil
	}
	// fmt.Println("testPath exist and check: ", testPath)
	tr, err := p.Parse(testPath, nil)
	if err != nil {
		if err == ErrEmptyFile {
			// Overwrite empty test files.
			return h, nil, nil
		}
		return nil, nil, fmt.Errorf("Parser.Parse test file: %v", err)
	}
	var testFuncs []string
	for _, fun := range tr.Funcs {
		// fmt.Println("fun.IsEcho 1: ", fun.IsEcho)
		if fun.IsEcho {
			h.Imports = append(h.Imports, &Import{
				Name: "",
				Path: `"github.com/codehand/cest/echo/mctx"`,
			})
		}
		// fun.Package = h.Package
		testFuncs = append(testFuncs, fun.Name)
	}
	tr.Header.Imports = append(tr.Header.Imports, h.Imports...)
	tr.Header.PkgName = h.PkgName
	h = tr.Header
	// fmt.Println("tr.Header:", tr.Header.Comments)
	return h, testFuncs, nil
}

func testableFuncs(h *Header, funcs []*Function, only, excl *regexp.Regexp, exp, out, sameDir bool, testFuncs []string) ([]*Function, *Header) {

	sort.Strings(testFuncs)
	var fs []*Function
	for _, f := range funcs {
		if isTestFunction(f, testFuncs) || isExcluded(f, excl) || isUnexported(f, exp) || !isIncluded(f, only) || isInvalid(f) {
			continue
		}
		h.Imports = append(h.Imports, &Import{
			Name: "",
			Path: `"github.com/stretchr/testify/assert"`,
		})

		if f.IsEcho {
			h.Imports = append(h.Imports, &Import{
				Name: "",
				Path: `"github.com/codehand/cest/echo/mctx"`,
			})
		}
		if !sameDir {
			f.Package = h.PkgName
		}
		// if out {

		// }
		fs = append(fs, f)
	}
	return fs, h
}

func isInvalid(f *Function) bool {
	if f.Name == "init" && f.IsNaked() {
		return true
	}
	return false
}

func isTestFunction(f *Function, testFuncs []string) bool {
	return len(testFuncs) > 0 && contains(testFuncs, f.TestName())
}

func isExcluded(f *Function, excl *regexp.Regexp) bool {
	return excl != nil && (excl.MatchString(f.Name) || excl.MatchString(f.FullName()))
}

func isUnexported(f *Function, exp bool) bool {
	if !f.IsExported {
		if exp {
			return false
		}
		return true
	}
	return false
}

func isIncluded(f *Function, only *regexp.Regexp) bool {
	return only == nil || only.MatchString(f.Name) || only.MatchString(f.FullName())
}

func contains(ss []string, s string) bool {
	if i := sort.SearchStrings(ss, s); i < len(ss) && ss[i] == s {
		return true
	}
	return false
}
