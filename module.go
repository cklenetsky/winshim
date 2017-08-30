package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type module struct {
	ModuleName     string                     // The base name of the output module
	SourceHeader   string                     // The C header file used to create the shim
	Functions      []parsedFunctionDefinition // The functions parsed from the C header file
	OutputDir      string                     // The directory to write the generated files
	OutputCFile    string                     // The name of the .c file to generate
	OutputGoFile   string                     // The base name of the .go files to generate
	OutputMakefile string                     // The name of the Makefile to generate
}

var funcMap = template.FuncMap{
	"tolower":      strings.ToLower,
	"convertslash": func(s string) string { return strings.Replace(s, "\\", "/", -1) },
	"filenamebase": func(s string) string {
		i := strings.LastIndex(s, ".")
		if i > 0 {
			return s[:i]
		}
		return s
	},
	"filepathfile": func(s string) string { _, f := filepath.Split(s); return f },
	"filepathpath": func(s string) string { p, _ := filepath.Split(s); return p },
}

func processTemplate(filename string, templatefile string, data interface{}) error {
	fileTemplate, err := template.New(templatefile).Funcs(funcMap).ParseFiles(templatefile)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	err = fileTemplate.ExecuteTemplate(f, templatefile, data)
	return err
}

func (m module) writeShimFile() error {
	return processTemplate(fmt.Sprintf("%s%c%s", m.OutputDir, os.PathSeparator, m.OutputCFile), "shimfile.template", m)
}

func (m module) writeMakefile() error {
	return processTemplate(fmt.Sprintf("%s%c%s", m.OutputDir, os.PathSeparator, m.OutputMakefile), "makefile.template", m)
}

func (m module) writeGofiles() error {
	type goFileInfo struct {
		module
		Windows bool
	}

	info := goFileInfo{
		module:  m,
		Windows: false,
	}

	err := processTemplate(fmt.Sprintf("%s%c%s.go", m.OutputDir, os.PathSeparator, strings.ToLower(m.OutputGoFile)), "gofile.template", info)
	if err != nil {
		return err
	}
	info.Windows = true
	return processTemplate(fmt.Sprintf("%s%c%s_windows.go", m.OutputDir, os.PathSeparator, strings.ToLower(m.OutputGoFile)), "gofile.template", info)
}
