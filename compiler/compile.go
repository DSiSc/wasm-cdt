package compiler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Compile(source, folderPath, libPath string) error {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("source file %s doesn't exist", source)
	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create output directory %s, as: %v", folderPath, err)
		}
	}
	var (
		baseInclude = fmt.Sprintf(`%s/misc/`, libPath)
		souceName   = strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
		wasmfile    = fmt.Sprintf(`%s/%s.wasm`, folderPath, souceName)
	)

	clang := clangCmd(source, baseInclude, wasmfile)
	if _, err := execCmd(clang); err != nil {
		return fmt.Errorf("failed to compile source file, as:%v", err)
	}
	return nil
}

func clangCmd(source, baseInclued, wasmfile string) *exec.Cmd {
	cmd := exec.Command(`clang`, `--target=wasm32`, `-O3`, `-nostdlib`, `-fno-builtin`, `-ffreestanding`, `-fno-threadsafe-statics`,
		`-nostdinc`, `-static`, `-fno-rtti`, `-fno-exceptions`, `-DBOOST_DISABLE_ASSERTS`, `-DBOOST_EXCEPTION_DISABLE`,
		`-Wl,--gc-sections,--merge-data-segments,-zstack-size=8192`, `-Wno-unused-command-line-argument`, `-s`, `-Xlinker`, `--export=invoke`,
		`-Xlinker`, `--no-entry`, `-Xlinker`, `--allow-undefined`,
		fmt.Sprintf(`-I%s/include/libcxx`, baseInclued),
		fmt.Sprintf(`-I%s/include/libc`, baseInclued),
		fmt.Sprintf(`-I%s/include/justitia`, baseInclued),
		fmt.Sprintf(`-L%s/lib`, baseInclued),
		`-lc`, `-lc++`, `-o`, wasmfile,
		fmt.Sprintf(`%s/src/justitia.c`, baseInclued), source)
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
