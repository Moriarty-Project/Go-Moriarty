package MoriartyCLI

import (
	"GoMoriarty/siteCrawler"
	"strconv"
	"time"

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
		Aliases: []string{"new"},
		Func:    cli.cliNewScanner,
		Help:    "sets up a new instance of the site scanner",
	})
	cmd.AddCmd(&ishell.Cmd{
		Name:    "Scan",
		Aliases: []string{"scan", "basic"},
		Func:    cli.cliScan,
		Help:    "sets up a new instance of the site scanner",
	})
	cli.MainShell.AddCmd(cmd)
}

func (cli *MoriartyCLI) cliNewScanner(c *ishell.Context) {
	// first, delete the old instance
	cli.scanData.crawler = nil

	// then, get info we'll need for the new one.
	c.Println("site endpoint to test: ")
	url := c.ReadLine()
	c.Println("method(blank for GET):")
	method := c.ReadLine()
	c.Println("path to save folder:")
	saveFolder := c.ReadLine()
	crawler, err := siteCrawler.NewBasicSiteCrawler(url, method, saveFolder)
	if err != nil {
		c.Println("error in making site crawler from given data.")
		c.Println(err)
		return
	}

	// we can assume everything worked
	c.Println("crawler initialized")
	cli.scanData.crawler = crawler
}

func (cli *MoriartyCLI) cliScan(c *ishell.Context) {
	for cli.scanData.crawler == nil {
		c.Println("you first need to have a scanner setup. Starting the process now!")
		cli.cliNewScanner(c)
	}
	// now lets scan!
	c.Println("Scanning site for the following ranges")
	c.Println("from:")
	start, err := strconv.ParseInt(c.ReadLine(), 0, 64)
	if err != nil {
		c.Println(" invalid start valued passed.")
		return
	}
	c.Println("to:")
	stop, err := strconv.ParseInt(c.ReadLine(), 0, 64)
	if err != nil {
		c.Println("invalid stop valued passed.")
		return
	}
	c.Println("max errors in a row(-1 to ignore):")
	maxErrs, err := strconv.ParseInt(c.ReadLine(), 0, 64)
	if err != nil {
		maxErrs = 250
	}
	iterator := siteCrawler.NewDefaultIterator(int(start), int(stop), 1)
	c.Println(" starting now!")
	foundAll := make(chan bool)
	go func() {
		defer close(foundAll)
		foundAll <- cli.scanData.crawler.GetAll(
			iterator,
			int(maxErrs),
		)
	}()
	c.ProgressBar().Final("finished scanning site!")
	c.ProgressBar().Start()
	progressedLast := start
	for {
		select {
		case found := <-foundAll:
			c.ProgressBar().Stop()
			if found {
				c.Println("found all without encountering error limit!")
			} else {
				c.Println("encountered the error limit before reaching the end")
			}
			return
		default:
			s, exists := iterator.Peak(1)
			if !exists {
				c.ProgressBar().Stop()
				c.Println("unknown error happened while waiting for scan to complete.")
				panic("unknown error happened while waiting for scan to complete")
			}
			progressedPoint, err := strconv.ParseInt(s, 0, 64)
			if err == nil && progressedPoint > progressedLast {
				// we can increment the progress bar
				c.ProgressBar().Progress(int((100 * (start - progressedLast - progressedPoint)) / stop))
				progressedLast = progressedPoint
			}
			<-time.After(time.Millisecond * 10)
		}
	}
}
