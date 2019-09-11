package main

import (
	"fmt"
	"github.com/DSiSc/craft/log"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func init() {
	log.Disable()
}

func NewApp(usage string) *cli.App {
	app := cli.NewApp()
	app.Version = "1.0"
	app.Usage = usage
	app.Commands = []cli.Command{
		CompileCommand,
		DeployCommand,
		AccountCommand,
	}
	return app
}

func main() {
	app := NewApp("Justitia WASM conract development tools")
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
