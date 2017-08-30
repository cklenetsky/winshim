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

// Signature returns the function's prototype signature, without the trailing semicolon
func (f parsedFunctionDefinition) Signature() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%s ", f.ReturnType))
	if f.Attribute != "" {
		b.WriteString(fmt.Sprintf("__%s ", f.Attribute))
	}
	b.WriteString(f.Name)
	b.WriteRune('(')
	b.WriteString(f.ParameterList())
	b.WriteRune(')')
	return b.String()
}

// String returns a string representation of this function definition
func (f parsedFunctionDefinition) String() string {
	return f.Signature()
}

// ParameterList returns this function's parameters in the form of
// void|<type> <var name>[, <type> <var name>, ...]
func (f parsedFunctionDefinition) ParameterList() string {
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
