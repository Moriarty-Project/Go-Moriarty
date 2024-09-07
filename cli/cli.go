package MoriartyCLI

// import (
// 	"GoMoriarty/moriarty"
// 	"fmt"

// 	"github.com/abiosoft/ishell/v2"
// )

// // the cli package for moriarty!
// type MoriartyCLI struct {
// 	MainShell     *ishell.Shell
// 	loadedUsers   []*moriarty.UserRecordings //the users we have loaded already
// 	unloadedUsers []string                   //the names of user profiles we unloaded.
// }

// func NewCLI() *MoriartyCLI {
// 	shell := ishell.New()
// 	cli := &MoriartyCLI{
// 		MainShell:     shell,
// 		loadedUsers:   []*moriarty.UserRecordings{},
// 		unloadedUsers: []string{},
// 	}
// 	cli.addUserMethods()

// 	shell.Println("Welcome to Moriarty!")

// 	shell.Start()
// 	return cli
// }
// func (cli *MoriartyCLI) addUserMethods() {
// 	cmd := &ishell.Cmd{
// 		Name:     "Users",
// 		Aliases:  []string{"users"},
// 		Help:     "add/edit/remove user profiles",
// 		LongHelp: "all controls related to user profiles!",
// 	}
// 	cmd.AddCmd(&ishell.Cmd{
// 		Name:    "Add",
// 		Aliases: []string{"add"},
// 		Func:    cli.cliAddUser,
// 		Help:    "add new user",
// 	})
// 	cmd.AddCmd(&ishell.Cmd{
// 		Name:    "Edit",
// 		Aliases: []string{"edit"},
// 		Func:    cli.cliEditUser,
// 		Help:    "select user profile to edit",
// 	})
// 	cmd.AddCmd(&ishell.Cmd{
// 		Name:    "Delete",
// 		Aliases: []string{"delete", "rm", "remove", "Remove"},
// 		Func:    cli.cliRemoveUser,
// 		Help:    "select users to remove",
// 	})
// 	cmd.AddCmd(&ishell.Cmd{
// 		Name:    "Load/Unload",
// 		Aliases: []string{"load", "Load", "unload", "Unload"},
// 		Func:    cli.cliLoadUnloadUser,
// 		Help:    "change users that are loaded vs stored.",
// 	})
// 	cli.MainShell.AddCmd(cmd)
// }
// func (cli *MoriartyCLI) cliAddUser(c *ishell.Context) {
// 	// TODO: add user functions
// 	// we'll need to get anything known about that user.
// }

// func (cli *MoriartyCLI) cliEditUser(c *ishell.Context) {
// 	// TODO: edit user functions
// 	//first, we need to let them select a user.
// 	cli.MainShell.Println("select user")
// 	cli.MainShell.Println("if an unloaded user is selected, they will be loaded to edit!")
// 	selectedUser := cli.SelectUser()

// 	if selectedUser == nil {
// 		c.Println("no user selected.")
// 		return
// 	}

// 	// edit the user
// 	options := make([]string, 0, 8)
// 	options = append(options, fmt.Sprintf("Name: %v", selectedUser.AccountName))
// 	options = append(options, "\t\tKnowns")
// 	options = append(options, "\t\tLikely")
// 	options = append(options, "\t\tPossible")
// 	options = append(options, "Back")
// 	selection := cli.MainShell.MultiChoice(options, "edit user")

// 	if selection == -1 {
// 		// they didn't select anything
// 		return
// 	}
// 	// USE: if the selection info is changed, this likely will need to be as well!
// 	switch selection {
// 	case 0:
// 		// they selected the username.
// 		cli.MainShell.Print("New Name:")
// 		cli.MainShell.ReadLineWithDefault(selectedUser.AccountName)
// 		// TODO: we might even want to check the file systems...
// 	case 1:
// 		// they want to edit the knowns
// 		cli.cliEditKnowns(c, selectedUser)
// 	case 2:
// 		// they want to edit the likely items
// 		cli.cliEditLikelys(c, selectedUser)
// 	case 3:
// 		// they want to edit the possible items
// 		cli.cliEditPossibles(c, selectedUser)
// 	default:
// 		return
// 	}
// }

// // edit knowns
// var editSelectedActions = []string{
// 	"Select upgradeable",
// 	"Select downloadable",
// 	"Select deletable",
// 	"Add username",
// 	"Add email",
// 	"back",
// }

// const badInputString = "nothing entered, this will be ignored"

// func selectToDelete(names []string, c *ishell.Context, user *moriarty.UserRecordings) {
// 	// delete
// 	accounts := c.Checklist(names, "select Items to remove", []int{})
// 	for _, i := range accounts {
// 		user.Delete(names[i])
// 	}
// }
// func (cli *MoriartyCLI) cliEditKnowns(c *ishell.Context, selectedUser *moriarty.UserRecordings) {
// 	selection := cli.MainShell.MultiChoice(editSelectedActions, "edit Knowns")
// 	names := append(selectedUser.KnownEmails, selectedUser.KnownUsernames...)
// 	switch selection {
// 	case 0:
// 		// do nothing, we cant upgrade...
// 		c.Println("cant upgrade from known...")
// 		c.Println("a less lazy dev might just remove the option")
// 	case 1:
// 		// downgrade
// 		accounts := c.Checklist(names, "select Items to downgrade to likely", []int{})
// 		for _, i := range accounts {
// 			selectedUser.Downgrade(names[i])
// 		}
// 	case 2:
// 		// delete
// 		selectToDelete(names, c, selectedUser)
// 	case 3:
// 		// new username
// 		c.Print("New Username:")
// 		ans := c.ReadLine()
// 		if ans != "" {
// 			selectedUser.AddKnownUsername(ans)
// 		} else {
// 			c.Println(badInputString)
// 		}
// 	case 4:
// 		// new email
// 		c.Print("New Email:")
// 		ans := c.ReadLine()
// 		if ans != "" {
// 			selectedUser.KnownEmails = append(selectedUser.KnownEmails, ans)
// 		} else {
// 			c.Println(badInputString)
// 		}
// 	}
// }

// func (cli *MoriartyCLI) cliEditLikelys(c *ishell.Context, selectedUser *moriarty.UserRecordings) {
// 	selection := cli.MainShell.MultiChoice(editSelectedActions, "edit Likely items")
// 	names := append(selectedUser.LikelyEmails, selectedUser.LikelyUsernames...)
// 	switch selection {
// 	case 0:
// 		// upgrade
// 		accounts := c.Checklist(names, "select Items to upgrade to known", []int{})
// 		for _, i := range accounts {
// 			selectedUser.Upgrade(names[i])
// 		}
// 	case 1:
// 		// downgrade
// 		accounts := c.Checklist(names, "select Items to downgrade to possible", []int{})
// 		for _, i := range accounts {
// 			selectedUser.Downgrade(names[i])
// 		}
// 	case 2:
// 		// delete
// 		selectToDelete(names, c, selectedUser)
// 	case 3:
// 		// new username
// 		c.Print("New Username:")
// 		ans := c.ReadLine()
// 		if ans != "" {
// 			// TODO: have this handle wildcards
// 			selectedUser.LikelyUsernames = append(selectedUser.LikelyUsernames, ans)
// 		} else {
// 			c.Println(badInputString)
// 		}
// 	case 4:
// 		// new email
// 		c.Print("New Email:")
// 		ans := c.ReadLine()
// 		if ans != "" {
// 			selectedUser.LikelyEmails = append(selectedUser.LikelyEmails, ans)
// 		} else {
// 			c.Println(badInputString)
// 		}
// 	}
// }
// func (cli *MoriartyCLI) cliEditPossibles(c *ishell.Context, selectedUser *moriarty.UserRecordings) {
// 	selection := cli.MainShell.MultiChoice(editSelectedActions, "edit Possible items")
// 	names := append(selectedUser.PossibleUsernames, selectedUser.PossibleEmails...)
// 	switch selection {
// 	case 0:
// 		// upgrade
// 		accounts := c.Checklist(names, "select Items to upgrade to likely", []int{})
// 		for _, i := range accounts {
// 			selectedUser.Upgrade(names[i])
// 		}
// 	case 1:
// 		// downgrade
// 		selectToDelete(names, c, selectedUser)
// 	case 2:
// 		// delete
// 		selectToDelete(names, c, selectedUser)
// 	case 3:
// 		// new username
// 		c.Print("New Username:")
// 		ans := c.ReadLine()
// 		if ans != "" {
// 			// TODO: have this handle wildcards
// 			selectedUser.PossibleUsernames = append(selectedUser.PossibleUsernames, ans)
// 		} else {
// 			c.Println(badInputString)
// 		}
// 	case 4:
// 		// new email
// 		c.Print("New Email:")
// 		ans := c.ReadLine()
// 		if ans != "" {
// 			selectedUser.PossibleEmails = append(selectedUser.PossibleEmails, ans)
// 		} else {
// 			c.Println(badInputString)
// 		}
// 	}
// }
// func (cli *MoriartyCLI) cliRemoveUser(c *ishell.Context) {
// 	// TODO: remove user functions
// }
// func (cli *MoriartyCLI) cliLoadUnloadUser(c *ishell.Context) {
// 	// TODO: loadUnload user functions
// }

// // returns the index of the selected user profile.
// func (cli *MoriartyCLI) SelectUser() (user *moriarty.UserRecordings) {
// 	selections := make([]string, 0, len(cli.loadedUsers)+len(cli.unloadedUsers)+2)
// 	for _, val := range cli.loadedUsers {
// 		// TODO: add a saved name value for user profiles
// 		selections = append(selections, val.KnownUsernames[0])
// 	}
// 	selections = append(selections, "VV-unloaded users-VV")
// 	selections = append(selections, cli.unloadedUsers...)

// 	ans := cli.MainShell.MultiChoice(selections, "select a user profile")
// 	if ans == -1 {
// 		return nil
// 	}
// 	if ans < len(cli.loadedUsers) {
// 		return cli.loadedUsers[ans]
// 	}
// 	ans -= len(cli.loadedUsers) + 1
// 	// TODO: load the user!
// 	return nil
// }
