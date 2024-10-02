package MoriartyCLI

import (
	"GoMoriarty/utils"

	"github.com/abiosoft/ishell/v2"
)

// the cli package for moriarty!
type MoriartyCLI struct {
	MainShell     *ishell.Shell
	loadedUsers   []*utils.UserRecordings //the users we have loaded already
	unloadedUsers []string                //the names of user profiles we unloaded.
	scanData      *scanInstance
}

func NewCLI() *MoriartyCLI {
	shell := ishell.New()
	cli := &MoriartyCLI{
		MainShell:     shell,
		loadedUsers:   []*utils.UserRecordings{},
		unloadedUsers: []string{},
		scanData:      &scanInstance{},
	}
	cli.addUserMethods()
	cli.addScanMethods()

	shell.Println("Welcome to Moriarty!")

	shell.Start()
	return cli
}
