package moriarty

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type dataFound struct {
	data map[string][]string
}

// automatically checks folder of data files, attempts to parse it, and get selected user's results.
type FolderBasedData struct {
	FolderPath string                //how to get to the folder, where the files are located.
	Name       string                //name to be associated with this data path
	loadedData map[string]*dataFound //we keep each loaded file saved for ease of searching.
}

// the highest level of scan possible. Goes over all items for their total data reports
func (fbd *FolderBasedData) GetData(searchCriteria string) *DataTestResults {
	// first, we'll create our found data object.
	found := NewDataTestResults(fbd.Name)
	// TODO: if any file matches the data, I want to add that whole file
	for key, fileData := range fbd.loadedData {
		if strings.Compare(key, searchCriteria) == 0 {
			// if the key is the same as our search criteria, then this likely has a known user? so we'll add that whole row I guess. Hope no one searches "email" or something like that...

			found.Add(key, fmt.Sprintf("%v", fileData.data))
			continue
		}
		// next, we need to try to check the values there.
		for valKey, valVal := range fileData.data {
			if strings.Contains(valKey, searchCriteria) {
				// then this value in this document is named after them.
				found.Add(key, append([]string{valKey}, valVal...)...)
				// continue if this was about them...
				continue
			}
			// next... each item needs to be checked...
			for _, subVal := range valVal {
				if strings.Contains(subVal, searchCriteria) {
					found.Add(key, valKey, subVal)
					found.Add(valKey, subVal)
				}
			}
		}
	}

	return found.NilIfEmpty()
}

// load all of the data from a given folder path into this memory.
func (fbd *FolderBasedData) LoadAllData(folderPath string, ignoredFiles ...string) error {
	// go to that folderPath, and give it a good bit of work to find it
	filePaths, err := fbd.getAllFileNames(folderPath)
	if err != nil {
		return err
	}
	// ok, we've either returned, or we have the dirs
	// go through all of the files in this directory, and open them. Ignoring sub directories.
	foundData := false
	// used to keep track of if we found any data
	for _, fileName := range filePaths {
		filePath := path.Join(folderPath, fileName)
		err = fbd.getDataFromFile(filePath)
		if err != nil {
			return err
		}
		foundData = true
	}
	if !foundData {
		// nothing was ever found in this folder.
		return fmt.Errorf("no valid files were found under path '%s'", folderPath)
	}
	return nil
}

// gets the absolute filepath to where ever you need. Attempts multiple things, but if none come back, returns the final error
// mostly is trying to find the correct file, returning an absolute path is purely by chance
func getAbsolutePath(to string) (string, error) {
	_, err := os.Stat(to)
	if err == nil {
		return to, nil
	}
	// that didn't work
	// next, try adding the cwd.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	_, err = os.Stat(path.Join(cwd, to))
	if err == nil {
		return path.Join(cwd, to), nil
	}
	// I guess we'll try merging them with ..'s until it fits?
	middle := "../"
	for i := 0; i < 5 && i < strings.Count(to, "/"); i++ {
		_, err = os.Stat(path.Join(cwd, middle, to))
		if err == nil {
			return path.Join(cwd, middle, to), nil
		}
		middle = middle + middle
	}
	return "", fmt.Errorf("could not find any path to %v", to)
}
func (fbd *FolderBasedData) getAllFileNames(folderPath string, ignoredFiles ...string) ([]string, error) {
	// go to that folderPath, and give it a good bit of work to find it
	folderPath, err := getAbsolutePath(folderPath)
	if err != nil {
		return nil, err
	}
	dirs, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	// quickly remove ignored files and directories before we go over the data
	filePaths := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		ignoredForReason := false
		for _, badName := range ignoredFiles {
			match, err := path.Match(badName, dir.Name())
			if err == nil && match ||
				strings.Compare(badName, dir.Name()) == 0 {
				ignoredForReason = true
				break
			}
		}
		if !ignoredForReason {
			// it's passed all tests, so we can count it as valid
			filePaths = append(filePaths, dir.Name())
		}
	}
	return filePaths, nil
}

/*
loads the data from that filepath into the loadedData.
current supported datatypes are json...
*/
func (fdb *FolderBasedData) getDataFromFile(filePath string) error {
	// check if the file ends with a format we can parse...
	if fdb.loadedData[path.Base(filePath)] != nil {
		return nil
	}
	suffix := path.Ext(filePath)

	// first we need to see if we can open the file
	file, err := fdb.getFileFrom(filePath)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return err
	}
	// now we need to parse all keywords from the file.
	dataRetrieved := map[string]interface{}{}
	dataBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	switch suffix {
	case ".json":
		// we need to parse it as a json object.
		json.Unmarshal(dataBytes, &dataRetrieved)
	default:
		return fmt.Errorf("unsupported format %v", suffix)
	}
	// we now have data retrieved from the file (hopefully).
	// lets parse it to strings from the interface, and work from there!
	found := &dataFound{
		data: map[string][]string{},
	}
	// lets parse it down!
	for key, val := range dataRetrieved {
		switch vt := val.(type) {
		case []string:
			found.data[key] = vt
		case []interface{}:
			found.data[key] = make([]string, 0, len(vt))
			for _, v := range vt {
				found.data[key] = append(found.data[key], fmt.Sprintf("%v", v))
			}
		default:
			// do our best to parse it with sprintf doing the legwork...
			found.data[key] = []string{fmt.Sprintf("%v", vt)}
		}
	}
	fdb.loadedData[path.Base(filePath)] = found
	return nil
}

// attempt to get the file. Tries multiple ways to get the file, and with a bit of luck, finds something!
func (fdb *FolderBasedData) getFileFrom(filePath string) (*os.File, error) {
	filePath, err := getAbsolutePath(filePath)
	if err != nil {
		return nil, err
	}
	return os.Open(filePath)
}

func (fbd *FolderBasedData) GetName() string { return fbd.Name }

// all folder based elements are assumed to be SFW
func (fbd *FolderBasedData) IsNsfw() bool { return false }
