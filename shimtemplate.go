package main

import "io"
import "bytes"
import "fmt"

const headers = `#include <windows.h>
`
const lockDefinition = `// These are available in Vista and newer
VOID WINAPI AcquireSRWLockShared(PSRWLOCK SRWLock);
VOID WINAPI ReleaseSRWLockShared(PSRWLOCK SRWLock);
VOID WINAPI AcquireSRWLockExclusive(PSRWLOCK SRWLock);
VOID WINAPI ReleaseSRWLockExclusive(PSRWLOCK SRWLock);

static SRWLOCK g_lock = SRWLOCK_INIT;
`

const libHandle = `// Per Microsoft:
// Simple reads and writes to properly aligned 64-bit variables are atomic on 
// 64-bit Windows.    
static HANDLE hLib = NULL;
`

const initFunctionStart = `int Init%s(const char *libname) {

	if (hLib != NULL)
		return 0;

	if (!libname || !libname[0])
		return -1;

	AcquireSRWLockExclusive(&g_lock);
    if ((hLib = LoadLibrary(libname)) == NULL) {
        ReleaseSRWLockExclusive(&g_lock);
        return -1;
    }		
`
const initFunctionEnd = `
	ReleaseSRWLockExclusive(&g_lock);
	return 0;
}
`
const cleanupFunction = `int Cleanup%s(void) {
    if (hLib != NULL) {
        AcquireSRWLockExclusive(&g_lock);
        FreeLibrary(hLib);
        hLib = NULL;
        ReleaseSRWLockExclusive(&g_lock);
    }
}
`

func writeNewline(output io.Writer) error {
	_, err := output.Write([]byte{'\n'})
	return err
}

func writeHeaderIncludes(output io.Writer) error {
	_, err := output.Write([]byte(headers))
	if err != nil {
		return err
	}
	err = writeNewline(output)
	if err != nil {
		return err
	}
	_, err = output.Write([]byte(lockDefinition))

	return err
}

func writeFunctionPointerVarTypes(functions []shimFunctionDefinition, output io.Writer) error {

	var b bytes.Buffer

	for _, f := range functions {
		b.WriteString(f.FunctionPointerTypedef())
		b.WriteRune('\n')
	}

	_, err := output.Write(b.Bytes())
	return err
}

func writeFunctionPointerVars(functions []shimFunctionDefinition, output io.Writer) error {

	var b bytes.Buffer

	for _, f := range functions {
		b.WriteString(fmt.Sprintf("static %s %s = NULL;\n", f.FunctionPointerTypedefName(), f.FunctionPointerName()))
	}
	_, err := output.Write(b.Bytes())
	return err
}

func writeGlobalHandleVar(output io.Writer) error {
	_, err := output.Write([]byte(libHandle))
	return err
}

func writeInitFunction(moduleName string, functions []shimFunctionDefinition, output io.Writer) error {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf(initFunctionStart, moduleName))

	for _, f := range functions {
		b.WriteString(fmt.Sprintf("    %s = (%s) GetProcAddress(hLib, \"%s\");", f.FunctionPointerName(), f.FunctionPointerTypedefName(), f.name))
		b.WriteRune('\n')
	}

	b.WriteString("\n    if (\n")
	for i, f := range functions {
		b.WriteString(fmt.Sprintf("        !%s", f.FunctionPointerName()))
		if i < len(functions)-1 {
			b.WriteString(" || \n")
		}
	}
	b.WriteString(` ) {
		
		FreeLibrary(hLib);
		hLib=NULL;
		ReleaseSRWLockExclusive(&g_lock);
		return -1;
	}`)

	b.WriteString(initFunctionEnd)

	_, err := output.Write(b.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func writeCleanupFunction(moduleName string, output io.Writer) error {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf(cleanupFunction, moduleName))

	_, err := output.Write(b.Bytes())
	return err
}

func writeShimFunctions(functions []shimFunctionDefinition, output io.Writer) error {
	var b bytes.Buffer

	for _, f := range functions {
		b.WriteString(f.ShimFunction("g_lock"))
		b.WriteRune('\n')
	}
	_, err := output.Write(b.Bytes())
	return err
}
