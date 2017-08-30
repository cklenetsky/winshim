package main

import (
	"io"
	"strings"
	"text/template"
)

type stubfileInfo struct {
	Windows    bool
	ModuleName string
}

func initGoStubTemplate() (*template.Template, error) {
	funcMap := template.FuncMap{
		"tolower": strings.ToLower,
	}
	return template.New("gofile.template").Funcs(funcMap).ParseFiles("gofile.template")
}

func writeWindowsGoStub(moduleName string, output io.Writer) error {

	stubfileTemplate, err := initGoStubTemplate()
	if err != nil {
		return err
	}
	stubStruct := stubfileInfo{
		Windows:    true,
		ModuleName: moduleName,
	}
	err = stubfileTemplate.ExecuteTemplate(output, "gofile.template", stubStruct)
	return err
}

func writeOtherGoStub(moduleName string, output io.Writer) error {
	stubfileTemplate, err := initGoStubTemplate()
	if err != nil {
		return err
	}
	stubStruct := stubfileInfo{
		Windows:    false,
		ModuleName: moduleName,
	}
	err = stubfileTemplate.Execute(output, stubStruct)
	return err
}
