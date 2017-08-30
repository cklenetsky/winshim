package main

import (
	"bytes"
	"fmt"
	"strings"
)

type parameterDefinition struct {
	name     string
	dataType string
}

type parsedFunctionDefinition struct {
	ReturnType string
	Name       string
	Attribute  string
	Parameters []parameterDefinition
}

func (f parsedFunctionDefinition) FunctionSignature() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%s ", f.ReturnType))
	if f.Attribute != "" {
		b.WriteString(fmt.Sprintf("__%s ", f.Attribute))
	}
	b.WriteString(f.Name)
	b.WriteRune('(')
	b.WriteString(f.FunctionParameters())
	b.WriteRune(')')
	return b.String()
}

func (f parsedFunctionDefinition) String() string {
	return f.FunctionSignature()
}

func (f parsedFunctionDefinition) FunctionParameters() string {
	var b bytes.Buffer

	if len(f.Parameters) == 0 {
		b.WriteString("void")
	} else {
		for i, p := range f.Parameters {
			b.WriteString(p.dataType)
			if !strings.HasSuffix(p.dataType, "*") {
				b.WriteRune(' ')
			}
			b.WriteString(p.name)
			if i < len(f.Parameters)-1 {
				b.WriteString(", ")
			}
		}
	}
	return b.String()
}
