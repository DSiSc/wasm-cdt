package main

import (
	"github.com/DSiSc/wasm-cdt/compiler"
	"gopkg.in/urfave/cli.v1"
)

var (
	CompileCommand = cli.Command{
		Action: CompileSource,
		Name:   "compile",
		Usage:  "Compile c/c++ file to wasm file",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "options, op",
				Usage: "compile `OPTIONS`",
			},
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
	options := ctx.String("options")
	sourceFile := ctx.String("file")
	outputPath := ctx.String("outpath")
	return compiler.Compile(options, sourceFile, outputPath)
}
