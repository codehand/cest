package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func RootFiles(srcPath string) ([]Path, string, error) {
	var rt = ""
	srcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return nil, rt, fmt.Errorf("filepath.Abs: %v", err)
	}
	var fi os.FileInfo
	if fi, err = os.Stat(srcPath); err != nil {
		return nil, rt, fmt.Errorf("os.Stat: %v", err)
	}
	if fi.IsDir() {
		// return dirFiles(srcPath)
	}

	// fmt.Println("isDir:", fi.IsDir())

	var files []Path
	var lop = true
	err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		// fmt.Println("zzz:", srcPath)

		if info.IsDir() || info.Name() == ".git" {
			return nil
		}
		if filepath.Ext(path) != ".go" || isHiddenFile(path) || strings.HasSuffix(string(path), "_test.go") {
			return nil
		}
		if lop {
			rt = path
			lop = false
		}
		files = append(files, Path(path))
		return nil
	})
	if err != nil {
		return nil, rt, err
	}
	return files, rt, nil
}

// Files returns all the Golang files for the given path. Ignores hidden files.
func Files(srcPath string) ([]Path, error) {
	srcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return nil, fmt.Errorf("filepath.Abs: %v", err)
	}
	var fi os.FileInfo
	if fi, err = os.Stat(srcPath); err != nil {
		return nil, fmt.Errorf("os.Stat: %v", err)
	}
	if fi.IsDir() {
		return dirFiles(srcPath)
	}
	return file(srcPath)
}

func dirFiles(srcPath string) ([]Path, error) {
	ps, err := filepath.Glob(path.Join(srcPath, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("filepath.Glob: %v", err)
	}
	var srcPaths []Path
	for _, p := range ps {
		src := Path(p)
		if isHiddenFile(p) || src.IsTestPath() {
			continue
		}
		srcPaths = append(srcPaths, src)
	}
	return srcPaths, nil
}

func file(srcPath string) ([]Path, error) {
	src := Path(srcPath)
	if filepath.Ext(srcPath) != ".go" || isHiddenFile(srcPath) {
		return nil, fmt.Errorf("no Go source files found at %v", srcPath)
	}
	return []Path{src}, nil
}

func isHiddenFile(path string) bool {
	return []rune(filepath.Base(path))[0] == '.'
}

func existOrCreateDir(src string) {
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		return
	}
	must(os.MkdirAll(src, 00755))
	printAction("green+h:black", "Created", ">> dir", src)
}

func existedDir(src string) bool {
	_, err := os.Stat(filepath.Dir(src))
	return err == nil
}
func ensureDir(src string) {
	srcN := filepath.Dir(src)
	// fmt.Println("src: ", srcN)
	if _, err := os.Stat(srcN); err != nil {
		err := os.MkdirAll(srcN, 00755)
		if err != nil {
			// panic(err)
			return
		}
		printAction("green+h:black", "Created", ">> file", src)
	}
}

func getParentDir(src string) string {
	return filepath.Dir(src)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// FileExists reports whether the named file exists as a boolean
func fileExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

// DirExists reports whether the dir exists as a boolean
func dirExists(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}

func fileMode(name string) int {
	if fileExists(name) {
		return 1
	}
	return 2
}
