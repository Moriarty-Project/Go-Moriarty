package moriarty

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
)

// automatically checks folder of data files, attempts to parse it, and get selected user's results.
type FolderBasedData struct {
	FolderPath string   //how to get to the folder, where the files are located.
	Name       string   //name to be associated with this data path
	files      []string //the file paths to each item
}

// if name is empty, it will attempt to create on from the last folder point in the folder path.
func NewFolderBasedData(folderPath string, name string, ignoredFiles ...string) (*FolderBasedData, error) {
	// check we can find our way to the folder
	folderPath, err := getAbsolutePath(folderPath)
	if err != nil {
		return nil, err
	}
	// next, check the folder path is valid.
	info, err := os.Stat(folderPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path does not lead to folder")
	}
	// check name has a value, if it doesn't, well give it one.
	if name == "" {
		name = info.Name()
	}
	fbd := &FolderBasedData{
		FolderPath: folderPath,
		Name:       name,
		files:      []string{},
	}
	return fbd, fbd.LoadAllData(folderPath, ignoredFiles...)
}

// the highest level of scan possible. Goes over all items for their total data reports
func (fbd *FolderBasedData) GetData(searchCriteria string) (*DataTestResults, error) {
	// first, we'll create our found data object.
	found := NewDataTestResults(searchCriteria)

	// then go through all of the files.
	for _, filePath := range fbd.files {
		file, err := fbd.getFileFrom(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if strings.Contains(filePath, searchCriteria) {
			//found them in the file name.
			found.Add("found in file name")
			found.Add("file name", file.Name())
			found.Add("file path", filePath)
			continue
		}
		fileReader := bufio.NewReader(file)
		if SearchBufferFor(fileReader, []byte(searchCriteria)) {
			found.Add("found in file contents")
			found.Add("file name", file.Name())
			found.Add("file path", filePath)
		}
	}

	return found.NilIfEmpty(), nil
}

// search function to search through an IO buffer for the following keyword, returns true if the keyword is found.
func SearchBufferFor(data *bufio.Reader, key []byte) bool {
	readData := make([]byte, len(key))
	for {
		// scan ahead until the first character of the key
		if _, err := data.ReadBytes(key[0]); err != nil {
			return false
		}
		data.UnreadByte()
		if _, err := data.Read(readData); err != nil {
			// we hit the end before it could begin
			return false
		}
		if bytes.Equal(key, readData) {
			// we found the key!
			return true
		}
		// reset the read data
		readData = make([]byte, len(key))
	}
}

// load all of the given files in this folder path into the files list
func (fbd *FolderBasedData) LoadAllData(folderPath string, ignoredFiles ...string) error {
	fileNames, err := fbd.getAllFileNames(folderPath, ignoredFiles...)
	if err != nil {
		return err
	}
	fbd.files = append(fbd.files, fileNames...)
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
			filePaths = append(filePaths, path.Join(folderPath, dir.Name()))
		}
	}
	return filePaths, nil
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
