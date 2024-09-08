package moriarty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const jsonFile = "resources/data.json"

func TestSiteCheckWithStatus(t *testing.T) {
	ubse := &StringBasedSiteElement{
		Name:                   "academia.edu",
		UrlHome:                "https://academia.edu",
		UrlUsername:            "https://independent.academia.edu/%s",
		UnclaimedIfResponseHas: "", //uses the status code, so doesn't matter.
		IsNSFW:                 false,
	}
	sut := NewSiteUserTester(ubse)
	assert.True(t, sut.TestSiteHas("Name")[0])
	assert.False(t, sut.TestSiteHas("reallyreallyrarename")[0])
}

func TestSiteCheckWithResponse(t *testing.T) {
	ubse := &StringBasedSiteElement{
		Name:                   "Apple Discussions",
		UrlHome:                "https://discussions.apple.com",
		UrlUsername:            "https://discussions.apple.com/profile/%s",
		UnclaimedIfResponseHas: "The page you tried was not found. You may have used an outdated link or may have typed the address (URL) incorrectly.",
		IsNSFW:                 false,
	}
	sut := NewSiteUserTester(ubse)
	assert.True(t, sut.TestSiteHas("helloWorld")[0])
	assert.False(t, sut.TestSiteHas("reallyasjnhnsujkdhnukdskjahjn")[0])
}

func TestLoadingElements(t *testing.T) {
	ubses, err := LoadAllSiteElements(jsonFile)
	assert.NoError(t, err)
	ubse := &StringBasedSiteElement{
		Name:                   "1337x",
		UrlHome:                "https://www.1337x.to/",
		UrlUsername:            "https://www.1337x.to/user/%s/",
		UnclaimedIfResponseHas: "\u003ctitle\u003eError something went wrong.\u003c/title\u003e",
		IsNSFW:                 false,
	}
	assert.Equal(t, ubse, ubses[0])
	siteTesters, err := LoadAllSiteUserTesters(jsonFile)
	assert.NoError(t, err)
	assert.Equal(t, ubse, siteTesters["1337x"].dataElement.(*StringBasedSiteElement))

	sut := siteTesters["1337x"]
	assert.True(t, sut.TestSiteHasAny("helloWorld"))
	assert.False(t, sut.TestSiteHasAny("reallyasjnhnsujkdhnukdskjahjn"))

}

func TestMoriartyUsability(t *testing.T) {
	arty, err := NewMoriarty(jsonFile)
	if err != nil {
		t.Fatal(err)
	}
	user := NewUserRecordings(
		"jonah",
	)
	arty.AssignNewUser(user)
	arty.trackingUser.AddNames("jpwilmsmeyer@gmail.com")
	assert.Empty(t, user.KnownFindings)
	assert.Empty(t, user.PossibleFindings)
	arty.TrackUserAcrossSites()

	// now check we found anything
	assert.NotEmpty(t, user.KnownFindings)
	// assert.NotEmpty(t, user.PossibleFindings)

}
