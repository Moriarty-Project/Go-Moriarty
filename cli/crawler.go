package MoriartyCLI

import (
	"GoMoriarty/siteCrawler"

	"github.com/abiosoft/ishell/v2"
)

type scanInstance struct {
	crawler *siteCrawler.SiteCrawler
}

func (cli *MoriartyCLI) addScanMethods() {
	cmd := &ishell.Cmd{
		Name:     "ScanSite",
		Aliases:  []string{"scan"},
		Help:     "scan a site",
		LongHelp: "test over api endpoint for basic iteration based scraping",
	}
	cmd.AddCmd(&ishell.Cmd{
		Name:    "New",
		Aliases: []string{},
		Func:    cli.cliNewScanner,
		Help:    "sets up a new instance of the site scanner",
	})
	cmd.AddCmd(&ishell.Cmd{
		Name:    "Scan",
		Aliases: []string{},
		Func:    cli.cliNewScanner,
		Help:    "sets up a new instance of the site scanner",
	})
	cli.MainShell.AddCmd(cmd)
}

func (cli *MoriartyCLI) cliNewScanner(c *ishell.Context) {
	// first, delete the old instance
	cli.scanData.crawler = nil

	// then, get info we'll need for the new one.
	c.Print("site endpoint to test: ")
	url := c.ReadLine()
	c.Print("\nmethod(blank for GET):")
	method := c.ReadLine()
	c.Print("\npath to save folder:")
	saveFolder := c.ReadLine()
	crawler, err := siteCrawler.NewBasicSiteCrawler(url, method, saveFolder)
	if err != nil {
		c.Println("\nerror in making site crawler from given data.")
		c.Println(err)
		return
	}

	// we can assume everything worked
	c.Println("\ncrawler initialized")
	cli.scanData.crawler = crawler
}
