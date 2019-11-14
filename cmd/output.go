package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/tools/imports"
)

type OptionsOutput struct {
	PrintInputs    bool
	Subtests       bool
	TemplateDir    string
	TemplateParams map[string]interface{}
}

func ProcessOutput(head *Header, funcs []*Function, opt *OptionsOutput) ([]byte, error) {
	if opt != nil && opt.TemplateDir != "" {
		err := LoadCustomTemplates(opt.TemplateDir)
		if err != nil {
			return nil, fmt.Errorf("loading custom templates: %v", err)
		}
	}

	tf, err := ioutil.TempFile("", "gotests_")
	if err != nil {
		return nil, fmt.Errorf("ioutil.TempFile: %v", err)
	}
	defer tf.Close()
	defer os.Remove(tf.Name())
	b := &bytes.Buffer{}
	if err := writeTests(b, head, funcs, opt); err != nil {
		return nil, err
	}

	// fmt.Println("b.Bytes():", string(b.Bytes()))
	out, err := imports.Process(tf.Name(), b.Bytes(), nil)
	if err != nil {
		return nil, fmt.Errorf("imports.Process: %v", err)
	}
	return out, nil
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func writeTests(w io.Writer, head *Header, funcs []*Function, opt *OptionsOutput) error {
	b := bufio.NewWriter(w)
	if err := renderHeader(b, head); err != nil {
		return fmt.Errorf("render.Header: %v", err)
	}
	for _, fun := range funcs {
		if err := renderTestFunction(b, head, fun, opt.PrintInputs, opt.Subtests, opt.TemplateParams); err != nil {
			return fmt.Errorf("render.TestFunction: %v", err)
		}
	}
	return b.Flush()
}
