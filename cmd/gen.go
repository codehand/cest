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
	srcFiles, err := Files(srcPath)
	if err != nil {
		return nil, fmt.Errorf("Files: %v", err)
	}
	files, err := Files(path.Dir(srcPath))
	if err != nil {
		return nil, fmt.Errorf("Files: %v", err)
	}
	// for _, item := range files {
	// 	fmt.Println(item.TestPath())
	// }
	// for _, item := range srcFiles {
	// 	fmt.Println(item.RealPath())
	// }
	if opt.Importer == nil || opt.Importer() == nil {
		opt.Importer = importer.Default
	}
	return parallelize(srcFiles, files, opt)
}

// result stores a generateTest result.
type result struct {
	gt  *GeneratedTest
	err error
}

// parallelize generates tests for the given source files concurrently.
func parallelize(srcFiles, files []Path, opt *Options) ([]*GeneratedTest, error) {
	var wg sync.WaitGroup
	rs := make(chan *result, len(srcFiles))
	for _, src := range srcFiles {
		wg.Add(1)
		// Worker
		go func(src Path) {
			defer wg.Done()
			r := &result{}
			r.gt, r.err = generateTest(src, files, opt)
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

func generateTest(src Path, files []Path, opt *Options) (*GeneratedTest, error) {
	p := &Parser{Importer: opt.Importer()}
	sr, err := p.Parse(string(src), files)
	fmt.Printf("check src: %v\n", string(src))
	fmt.Printf("check files: %v\n", files)
	fmt.Printf("check sr: %v\n", sr.Funcs)
	// find func
	for _, item := range sr.Funcs {
		fmt.Printf("func: %v\n", item)
	}

	if err != nil {
		return nil, fmt.Errorf("Parser.Parse source file: %v", err)
	}
	h := sr.Header
	h.Code = nil // Code is only needed from parsed test files.
	testPath := Path(src).TestPath()
	fmt.Printf("func: 1 %s\n", testPath)
	h, tf, err := parseTestFile(p, testPath, h)
	if err != nil {
		return nil, err
	}

	fmt.Printf("lst cmt: %v\n", h.Comments) // tamnt add comment
	funcs := testableFuncs(sr.Funcs, opt.Only, opt.Exclude, opt.Exported, tf)
	if len(funcs) == 0 {
		return nil, nil
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
	fmt.Printf("func: 2\n")
	return &GeneratedTest{
		Path:      testPath,
		Functions: funcs,
		Output:    b,
	}, nil
}

func parseTestFile(p *Parser, testPath string, h *Header) (*Header, []string, error) {
	if !IsFileExist(testPath) {
		fmt.Println("return if not exist")
		return h, nil, nil
	}
	fmt.Println("testPath exist and check: ", testPath)
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
		testFuncs = append(testFuncs, fun.Name)
	}
	tr.Header.Imports = append(tr.Header.Imports, h.Imports...)
	h = tr.Header
	fmt.Println("tr.Header:", tr.Header.Comments)
	return h, testFuncs, nil
}

func testableFuncs(funcs []*Function, only, excl *regexp.Regexp, exp bool, testFuncs []string) []*Function {
	sort.Strings(testFuncs)
	var fs []*Function
	for _, f := range funcs {
		if isTestFunction(f, testFuncs) || isExcluded(f, excl) || isUnexported(f, exp) || !isIncluded(f, only) || isInvalid(f) {
			continue
		}
		fs = append(fs, f)
	}
	return fs
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
	return exp && !f.IsExported
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
