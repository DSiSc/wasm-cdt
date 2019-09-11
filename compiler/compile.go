package compiler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"
)

const (
	WINDOWS = "windows"
	LINUX   = "linux"
)

func Compile(options, source, folderPath string) error {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("source file %s doesn't exist", source)
	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create output directory %s, as: %v", folderPath, err)
		}
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current path, as:%v", err)
	}

	var cmdPath string
	if LINUX == runtime.GOOS {
		cmdPath = currentDir + "/tools/linux/"
	} else if WINDOWS == runtime.GOOS {
		cmdPath = currentDir + "/tools/windows/"
	} else {
		return fmt.Errorf("unsupported os type %s", runtime.GOOS)
	}

	var (
		stdOut      string
		baseInclude = currentDir + "/misc/include"
		souceName   = strings.TrimSuffix(path.Base(source), path.Ext(source))
		bcFile      = fmt.Sprintf(`%s/%s.bc`, folderPath, souceName)
		sfile       = fmt.Sprintf(`%s/%s.s`, folderPath, souceName)
		watfile     = fmt.Sprintf(`%s/%s.wast`, folderPath, souceName)
		wasmfile    = fmt.Sprintf(`%s/%s.wasm`, folderPath, souceName)
	)

	clang := clangCmd(cmdPath, options, source, bcFile, baseInclude)
	if _, err := execCmd(clang); err != nil {
		return fmt.Errorf("failed to compile source file to bcfile, as:%v", err)
	}
	defer os.Remove(bcFile)

	llc := llcCmd(cmdPath, sfile, bcFile)
	if _, err := execCmd(llc); err != nil {
		return fmt.Errorf("failed to assemble bcfile to sfile, as:%v", err)
	}
	replaceSfile(sfile)
	defer os.Remove(sfile)

	s2wasm := s2wasmCmd(cmdPath, sfile)
	if stdOut, err = execCmd(s2wasm); err != nil {
		return fmt.Errorf("failed to compile sfile to wat file, as:%v", err)
	}
	if f, err := os.Create(watfile); err != nil {
		return fmt.Errorf("failed to create wat file, as:%v", err)
	} else if _, err := f.WriteString(stdOut); err != nil {
		return fmt.Errorf("failed to write content to wat file, as: %v", err)
	}

	wat2wasm := wat2wasmCmd(cmdPath, watfile, wasmfile)
	if _, err := execCmd(wat2wasm); err != nil {
		return fmt.Errorf("failed to convert wat file to wasm file, as:%v", err)
	}
	return nil
}

func clangCmd(cmdPath, options, source, bcfile, baseInclued string) *exec.Cmd {
	cmdStr := cmdPath + "llvm/bin/clang"
	cmd := exec.Command(cmdStr, `-emit-llvm`, `-nostdinc`, fmt.Sprintf(`-I%s/compat`, baseInclued),
		fmt.Sprintf(`-I%s/libcxx`, baseInclued), fmt.Sprintf(`-I%s/libc`, baseInclued), `-U__APPLE__`,
		`-D__EMSCRIPTEN__`, `--target=wasm32`, options, source, `-c`, `-o`, bcfile)
	return cmd
}

func replaceSfile(sfile string) {
	r, _ := regexp.Compile(`(\t\.import_global[^\n]*)\n\t\.size[^\n]*`)
	r.ReplaceAllString(sfile, `$1`)
}

func llcCmd(cmdPath, sfile, bcfile string) *exec.Cmd {
	cmd := exec.Command(cmdPath+"llvm/bin/llc", `-asm-verbose=false`, `-o`, sfile, bcfile)
	return cmd
}

func s2wasmCmd(cmdPath, sfile string) *exec.Cmd {
	cmd := exec.Command(cmdPath+"binaryen/s2wasm", sfile)
	return cmd
}

func wat2wasmCmd(cmdPath, watfile, wasmfile string) *exec.Cmd {
	cmd := exec.Command(cmdPath+"binaryen/wat2wasm", watfile, `-o`, wasmfile)
	return cmd
}

func execCmd(cmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(stderr.String())
	}
	return stdout.String(), nil
}
