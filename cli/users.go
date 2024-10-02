package MoriartyCLI

// the user based parts of the CLI
import (
	"GoMoriarty/utils"
	"strings"

	"github.com/abiosoft/ishell/v2"
)

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
	c.Println("Now Adding a new user!")

	// TODO: get the user file path, check the path is actually valid and wont break things.
	c.Print("What should this record be called?:")
	name := c.ReadLine()
	if !strings.HasSuffix(name, ".json") {
		name += ".json"
	}
	newUser := utils.NewUserRecordings(name)
	c.Println("are there any known account names? (leave blank for no)")
	for val := c.ReadLine(); val != ""; val = c.ReadLine() {
		newUser.AddNames(val)
	}

	// assuming nothing's caused us to cancel by now.
	cli.loadedUsers = append(cli.loadedUsers, newUser)
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
