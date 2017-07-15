package main

import "io"
import "bytes"
import "fmt"

func writeNewline(output io.Writer) error {
	_, err := output.Write([]byte{'\n'})
	return err
}

func writeHeaderIncludes(output io.Writer) error {
	_, err := output.Write([]byte("#include <windows.h>\n"))
	return err
}

func writeFunctionPointerVarTypes(functions []functionDefinition, output io.Writer) error {

	var b bytes.Buffer

	for _, f := range functions {

		b.WriteString(f.FunctionPointerTypedef())
		b.WriteRune('\n')
	}

	_, err := output.Write(b.Bytes())
	return err
}

func writeFunctionPointerVars(functions []functionDefinition, output io.Writer) error {

	var b bytes.Buffer

	for _, f := range functions {

		b.WriteString(fmt.Sprintf("static %s fp%s = NULL;\n", f.FunctionPointerTypedefName(), f.name))
	}
	_, err := output.Write(b.Bytes())
	return err
}

func writeGlobalHandleVar(output io.Writer) error {
	_, err := output.Write([]byte("static HANDLE hLib = NULL;\n"))
	return err
}

func writeInitFunction(moduleName string, dllFilename string, functions []functionDefinition, output io.Writer) error {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf(`int Init%s(void) {

	hLib = LoadLibrary("%s");
	if (!hLib) {
		return 1;
	}
	
`, moduleName, dllFilename))

	for _, f := range functions {
		b.WriteString(fmt.Sprintf("    fp%s = (%s) GetProcAddress(hLib, \"%s\");", f.name, f.FunctionPointerTypedefName(), f.name))
		b.WriteRune('\n')
	}

	b.WriteString("\n    if (\n")
	for i, f := range functions {
		b.WriteString(fmt.Sprintf("        !fp%s", f.name))
		if i < len(functions)-1 {
			b.WriteString(" || \n")
		}
	}
	b.WriteString(" ) {\n\n        FreeLibrary(hLib);\n        hLib=NULL;\n        return 1;\n    }\n")

	b.WriteString(`
	return 0;
}`)
	_, err := output.Write(b.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func writeCleanupFunction(moduleName string, output io.Writer) error {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf(`int Cleanup%s(void) {
	if (hLib) {
		return FreeLibrary(hLib)?0:1;
	}
}
`, moduleName))

	_, err := output.Write(b.Bytes())
	return err
}

func writeShimFunctions(functions []functionDefinition, output io.Writer) error {
	var b bytes.Buffer

	for _, f := range functions {
		b.WriteRune('\n')
		b.WriteString(f.ShimFunction())
		b.WriteRune('\n')
	}
	_, err := output.Write(b.Bytes())
	return err
}
