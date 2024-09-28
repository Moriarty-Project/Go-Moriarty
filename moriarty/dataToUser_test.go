package moriarty

import (
	"GoMoriarty/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSiteHasUser(t *testing.T) {
	sbse := &StringBasedSiteElement{
		Name:                   "GitHub",
		UrlHome:                "https://www.github.com/",
		UrlUsername:            "https://www.github.com/%s",
		UnclaimedIfResponseHas: "",
		IsNSFW:                 false,
	}
	sut := NewDataToUserTester(sbse)
	user := utils.NewUserRecordings("test myself")
	user.AddNames("PyroHedgehog")
	newKnown, newPossible, err := sut.TestSourceWith(user)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, newKnown)
	assert.Empty(t, newPossible)
}

func TestFolderSearcher(t *testing.T) {
	fbd, err := NewFolderBasedData("moriarty/resources/testingData/", "local folder")
	if err != nil {
		t.Fatal(err)
	}
	dt := NewDataToUserTester(fbd)
	assert.False(t, dt.IsNSFW())
	assert.Equal(t, dt.GetSourceName(), "local folder")
	results := dt.GetSourceResults("testEmail@test.com")
	assert.NotNil(t, results)

}
