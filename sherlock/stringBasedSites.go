package sherlock

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func LoadAllSiteElements(filePath string) ([]*StringBasedSiteElement, error) {
	// first, check if the path leads to a large data.json file...
	if !strings.HasSuffix(filePath, ".json") {
		// at some point we might want to deal with this type, eg, multiple files within a folder, but for now, we just throw an error
		return nil, fmt.Errorf("path is not to a JSON file")
	}
	if !path.IsAbs(filePath) {
		// if it's not an absolute path, add the CWD.
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		filePath = path.Join(cwd, filePath)
	}
	// now we try to find the file.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fileData, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		return nil, err
	}
	ubses := make([]*StringBasedSiteElement, 0, 100)
	err = json.Unmarshal(fileData, &ubses)
	if err != nil {
		return ubses, err
	}
	// I've split this up so we can check if we want to do any further sanitization here!

	return ubses, nil
}

// this is the base way that sites are saved.
type StringBasedSiteElement struct {
	Name                   string //the name of the site
	UrlHome                string //the home URL for the site
	UrlUsername            string //the url that has a way to format in the username
	UnclaimedIfResponseHas string //if the response text has this, then it is unclaimed.
	IsNSFW                 bool   //is this site NSFW
}
