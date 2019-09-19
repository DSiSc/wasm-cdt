package main

import (
	"fmt"
	"github.com/DSiSc/wasm-cdt/compiler"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
)

var (
	CompileCommand = cli.Command{
		Action: CompileSource,
		Name:   "compile",
		Usage:  "Compile c/c++ file to wasm file",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "c/c++ source code `FILE`",
			},
			cli.StringFlag{
				Name:  "outpath, d",
				Usage: "wasm file output `PATH`",
			},
		}}
)

func CompileSource(ctx *cli.Context) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("failed to get lib path, as: %v", err)
	}
	sourceFile := ctx.String("file")
	outputPath := ctx.String("outpath")
	return compiler.Compile(sourceFile, outputPath, dir)
}
