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

type shimFunctionDefinition struct {
	returnType string
	name       string
	attribute  string
	parameters []parameterDefinition
}

const functionTypePrefix = "f_"
const functionPointerVarPrefix = "fp"

func (f shimFunctionDefinition) FunctionSignature() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%s ", f.returnType))
	if f.attribute != "" {
		b.WriteString(fmt.Sprintf("__%s ", f.attribute))
	}
	b.WriteString(f.name)
	b.WriteRune('(')
	if len(f.parameters) == 0 {
		b.WriteString("void")
	} else {
		for i, p := range f.parameters {
			b.WriteString(p.dataType)
			if !strings.HasSuffix(p.dataType, "*") {
				b.WriteRune(' ')
			}
			b.WriteString(p.name)
			if i < len(f.parameters)-1 {
				b.WriteString(", ")
			}
		}
	}
	b.WriteRune(')')
	return b.String()
}

func (f shimFunctionDefinition) String() string {
	return f.FunctionSignature()
}

func (f shimFunctionDefinition) FunctionPointerName() string {
	return fmt.Sprintf("%s%s", functionPointerVarPrefix, f.name)
}

func (f shimFunctionDefinition) FunctionPointerTypedefName() string {
	return fmt.Sprintf("%s%s", functionTypePrefix, f.name)
}

func (f shimFunctionDefinition) FunctionPointerTypedef() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("typedef %s (", f.returnType))
	if f.attribute != "" {
		b.WriteString(fmt.Sprintf("__%s ", f.attribute))
	}
	b.WriteString(fmt.Sprintf("*%s)(", f.FunctionPointerTypedefName()))
	if len(f.parameters) == 0 {
		b.WriteString("void")
	} else {
		for i, p := range f.parameters {
			b.WriteString(p.dataType)
			if !strings.HasSuffix(p.dataType, "*") {
				b.WriteRune(' ')
			}
			b.WriteString(p.name)
			if i < len(f.parameters)-1 {
				b.WriteString(", ")
			}
		}
	}
	b.WriteString(");")
	return b.String()
}

func (f shimFunctionDefinition) ShimFunction(lockVarName string) string {
	var b bytes.Buffer

	b.WriteString(f.String())
	b.WriteString(" {\n")
	if f.returnType != "void" {
		b.WriteString(fmt.Sprintf("    %s res;\n", f.returnType))
	}
	b.WriteString(fmt.Sprintf("    AcquireSRWLockShared(&%s);\n", lockVarName))
	if f.returnType != "void" {
		b.WriteString("    res = ")
	} else {
		b.WriteString("    ")
	}
	b.WriteString(fmt.Sprintf("%s(", f.FunctionPointerName()))
	if len(f.parameters) > 0 {
		for i, p := range f.parameters {
			b.WriteString(p.name)
			if i < len(f.parameters)-1 {
				b.WriteString(", ")
			}
		}
	}
	b.WriteString(");\n")
	b.WriteString(fmt.Sprintf("    ReleaseSRWLockShared(&%s);\n", lockVarName))
	if f.returnType != "void" {
		b.WriteString("    retun res;\n")
	}

	b.WriteString("}\n")
	return b.String()
}
