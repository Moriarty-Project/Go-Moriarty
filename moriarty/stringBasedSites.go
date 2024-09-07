package moriarty

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
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
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	ubses := make([]*StringBasedSiteElement, 0, 100)
	err = json.Unmarshal(file, &ubses)
	if err != nil {
		return ubses, err
	}
	// I've split this up so we can check if we want to do any further sanitization here!

	return ubses, nil
}

func WriteAllSiteElements(filePath string, elements []*StringBasedSiteElement) error {
	// first, check if the path leads to a large data.json file...
	if !strings.HasSuffix(filePath, ".json") {
		// at some point we might want to deal with this type, eg, multiple files within a folder, but for now, we just throw an error
		return fmt.Errorf("path is not to a JSON file")
	}
	if !path.IsAbs(filePath) {
		// if it's not an absolute path, add the CWD.
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		filePath = path.Join(cwd, filePath)
	}

	// now we try to find the file.
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	// write to the file!
	data, err := json.MarshalIndent(elements, "", "  ")
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	file.Close()
	return nil
}
func AddSiteElement(filePath string, newItems ...*StringBasedSiteElement) error {
	data, err := LoadAllSiteElements(filePath)
	if err != nil {
		return err
	}
	data = append(data, newItems...)
	// now sort them by name
	sort.Slice(data, func(i, j int) bool {
		return strings.Compare(data[i].Name, data[j].Name) == -1
	})
	return WriteAllSiteElements(filePath, data)
}

// this is the base way that sites are saved.
type StringBasedSiteElement struct {
	Name                   string //the name of the site
	UrlHome                string //the home URL for the site
	UrlUsername            string //the url that has a way to format in the username
	UnclaimedIfResponseHas string //if the response text has this, then it is unclaimed.
	IsNSFW                 bool   //is this site NSFW
}

// returns nil if nothing is found!
func (sbse *StringBasedSiteElement) GetData(searchCriteria string) (*DataTestResults, error) {
	if !sbse.TestHas(searchCriteria) {
		// we dont have it, so return nil without an error
		return nil, nil
	}
	// we have it!
	dtr := NewDataTestResults(sbse.Name)
	dtr.Add("Valid Name", searchCriteria)
	dtr.Add("Home URL", sbse.UrlHome)
	dtr.Add("Found At", fmt.Sprintf(sbse.UrlUsername, searchCriteria))
	return dtr, nil
}
func (sbse *StringBasedSiteElement) GetName() string {
	return sbse.Name
}
func (sbse *StringBasedSiteElement) IsNsfw() bool {
	return sbse.IsNSFW
}

// Test if the site has any positive response from this name
func (sbse *StringBasedSiteElement) TestHas(searchCriteria string) bool {
	response, err := http.Get(fmt.Sprintf(sbse.UrlUsername, searchCriteria))
	if err != nil {
		return false
	}
	// now we need to parse the response to see if it was found.

	// assume that a 404 is always a fail. So we check that it's all 200's
	if response.StatusCode != 200 {
		return false
	}
	// check if it was found properly. We do this by checking if it doesn't have the false case string in its body.
	found := true
	if sbse.UnclaimedIfResponseHas != "" {
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return false
		}
		found = !strings.Contains(string(responseData), sbse.UnclaimedIfResponseHas)
	}
	return found
}
