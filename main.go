package main

import MoriartyCLI "GoMoriarty/cli"

func main() {
	cli := MoriartyCLI.NewCLI()
	cli.MainShell.Run()
}
