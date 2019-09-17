package compiler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Compile(source, folderPath string) error {
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

	var (
		baseInclude = currentDir + "/misc/include"
		souceName   = strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
		iRfile      = fmt.Sprintf(`%s/%s.o`, folderPath, souceName)
		wasmfile    = fmt.Sprintf(`%s/%s.wasm`, folderPath, souceName)
	)

	clang := clangCmd(source, iRfile, baseInclude)
	if _, err := execCmd(clang); err != nil {
		return fmt.Errorf("failed to compile source file, as:%v", err)
	}
	defer os.Remove(iRfile)

	lld := wasmLdCmd(iRfile, wasmfile)
	if _, err := execCmd(lld); err != nil {
		return fmt.Errorf("failed to generate wasm file, as:%v", err)
	}
	return nil
}

func clangCmd(source, iRfile, baseInclued string) *exec.Cmd {
	cmd := exec.Command(`clang`, fmt.Sprintf(`-I%s/compat`, baseInclued),
		fmt.Sprintf(`-I%s/libcxx`, baseInclued), fmt.Sprintf(`-I%s/libc`, baseInclued), `-c`, `-O3`, `--target=wasm32`, `-o`, iRfile, source)
	return cmd
}

func wasmLdCmd(iRfile, wasmfile string) *exec.Cmd {
	cmd := exec.Command(`wasm-ld`, `--no-entry`, `--strip-all`, `--allow-undefined`, `--export-all`, iRfile, `-o`, wasmfile)
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
