package moriarty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderBasedData(t *testing.T) {
	fbd := &FolderBasedData{
		FolderPath: "",
		Name:       "TestFolder",
	}
	// test while it's empty
	rslt, err := fbd.GetData("test")
	if err != nil {
		t.Fatal(err)
	}
	if rslt != nil {
		t.Logf("expected nil, got %v", rslt)
		t.Fail()
	}
	// load some dummy data
	err = fbd.LoadAllData("moriarty/resources/testingData/")
	assert.NoError(t, err)

	// now we can test searching real data
	rslt, err = fbd.GetData("Firstname")
	if err != nil {
		t.Fatal(err)
	}
	// just check we get something out of this
	assert.NotNil(t, rslt)
	rslt, err = fbd.GetData("UserDataDemo")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "UserDataDemo", rslt.Name)
	assert.True(t, rslt.InfoFound["found in file name"] != nil)

}
