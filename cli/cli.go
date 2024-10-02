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
}

func NewCLI() *MoriartyCLI {
	shell := ishell.New()
	cli := &MoriartyCLI{
		MainShell:     shell,
		loadedUsers:   []*utils.UserRecordings{},
		unloadedUsers: []string{},
	}
	cli.addUserMethods()

	shell.Println("Welcome to Moriarty!")

	shell.Start()
	return cli
}
func (cli *MoriartyCLI) addUserMethods() {
	cmd := &ishell.Cmd{
		Name:     "Users",
		Aliases:  []string{"users"},
		Help:     "add/edit/remove user profiles",
		LongHelp: "all controls related to user profiles!",
	}
	cmd.AddCmd(&ishell.Cmd{
		Name:    "Add",
		Aliases: []string{"add"},
		Func:    cli.cliAddUser,
		Help:    "add new user",
	})
	cmd.AddCmd(&ishell.Cmd{
		Name:    "Edit",
		Aliases: []string{"edit"},
		Func:    cli.cliEditUser,
		Help:    "select user profile to edit",
	})
	cmd.AddCmd(&ishell.Cmd{
		Name:    "Load",
		Aliases: []string{"load", "Load"},
		Func:    cli.cliLoadUser,
		Help:    "change users that are loaded vs stored.",
	})
	cli.MainShell.AddCmd(cmd)
}
func (cli *MoriartyCLI) cliAddUser(c *ishell.Context) {
	// TODO: add user functions
	// a simple series of questions to get the users info, then we'll go from there.
	// get the name for the user, where to save them, and any known info on them.
}

func (cli *MoriartyCLI) cliEditUser(c *ishell.Context) {
	cli.MainShell.Println("just edit the users JSON file before loading them.")
}

func (cli *MoriartyCLI) cliLoadUser(c *ishell.Context) {
	// TODO: load user from json
}

// returns the index of the selected user profile.
func (cli *MoriartyCLI) SelectUser() (user *utils.UserRecordings) {
	selections := make([]string, 0, len(cli.loadedUsers)+len(cli.unloadedUsers)+2)
	for _, val := range cli.loadedUsers {
		selections = append(selections, val.AccountName)
	}
	selections = append(selections, "VV-unloaded users-VV")
	selections = append(selections, cli.unloadedUsers...)

	ans := cli.MainShell.MultiChoice(selections, "select a user profile")
	if ans == -1 {
		return nil
	}
	if ans < len(cli.loadedUsers) {
		return cli.loadedUsers[ans]
	}
	ans -= len(cli.loadedUsers) + 1
	// TODO: load the user!
	return nil
}
