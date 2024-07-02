package moriarty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderBasedData(t *testing.T) {
	fbd := &FolderBasedData{
		FolderPath: "",
		Name:       "TestFolder",
		loadedData: map[string]*dataFound{},
	}
	// test while it's empty
	rslt := fbd.GetData("test")
	if rslt != nil {
		t.Logf("expected nil, got %v", rslt)
		t.Fail()
	}
	// load some dummy data
	err := fbd.LoadAllData("moriarty/resources/testingData/")
	assert.NoError(t, err)

	// now we can test searching real data
	rslt = fbd.GetData("Firstname")

	// just check we get something out of this
	assert.NotNil(t, rslt)

}
